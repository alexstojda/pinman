ARG RAILWAY_STATIC_URL="localhost:8080"
ARG PROTOCOL="http://"
ARG REACT_APP_API_HOST="$PROTOCOL$RAILWAY_STATIC_URL/api"

#######
# Generate Axios Typescript API Client
#######
FROM openapitools/openapi-generator-cli:latest AS node-gen

COPY api/.openapi-generator-ignore /out/
COPY api/openapi.yaml /api/

RUN /usr/local/bin/docker-entrypoint.sh generate \
-i /api/openapi.yaml \
-g typescript-axios \
-o /out

#######
# Build SPA
#######
FROM node:18-alpine AS node-dev

RUN mkdir -p /app
WORKDIR /app

COPY web/app/package.json /app
RUN yarn

COPY web/app/ /app
COPY --from=node-gen /out/ /app/src/api/generated

ARG REACT_APP_API_HOST
ENV REACT_APP_API_HOST $REACT_APP_API_HOST
RUN yarn build

ENTRYPOINT ["yarn"]

#######
# Golang shared base image
#######
FROM golang:1.19-alpine AS go-base

#######
# Generate Golang API Server
#######
FROM go-base AS go-gen

RUN go install "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.0"
RUN mkdir -p /out

COPY api /api
RUN /go/bin/oapi-codegen -config /api/oapi-codegen.config.yaml /api/openapi.yaml > /out/openapi.gen.go

#######
# Build API Server
#######
FROM go-base AS go-dev

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN openssh-client ca-certificates && update-ca-certificates 2>/dev/null || true
RUN apk add --no-cache git make

ENV HOME=/home/golang
WORKDIR /app
RUN adduser -h $HOME -D -u 1000 -G root golang && \
  chown golang:root /app && \
  chmod g=u /app $HOME
USER golang:root

COPY --chown=golang:root go.mod go.sum Makefile ./

RUN make mod

COPY --chown=golang:root cmd ./cmd
COPY --chown=golang:root internal ./internal
COPY --from=go-gen /out/ internal/app/generated/
RUN go build -v -o pinman ./cmd/pinman/

ENTRYPOINT ["make"]
CMD ["test"]

#######
# Final production image with all assets
#######
FROM scratch AS prod

COPY --from=go-dev /etc/passwd /etc/group  /etc/
COPY --from=go-dev /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Kube crashes if there isn't a tmp directory to write error logs to
COPY --from=go-dev --chown=golang:root /tmp /tmp
COPY --from=go-dev --chown=golang:root /app/pinman /app/
COPY --from=node-dev --chown=golang:root /app/build /app/html

COPY --chown=golang:root .env /app/

USER golang:root
EXPOSE 8080

ENV SPA_PATH=/app/html
WORKDIR /app
ENTRYPOINT ["./pinman"]
