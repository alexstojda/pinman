FROM node:18-alpine AS node-dev

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY package.json /usr/src/app
RUN yarn
COPY . /usr/src/app
RUN yarn build

ENTRYPOINT ["npm"]

FROM golang:1.19-alpine AS go-dev

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

COPY --chown=golang:root . ./
RUN go build -v -o pinman ./cmd/pinman/
RUN go build -v -o migrate ./cmd/migrate/

ENTRYPOINT ["make"]
#CMD ["test"]

###

FROM scratch AS prod

COPY --from=go-dev /etc/passwd /etc/group  /etc/
COPY --from=go-dev /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Kube crashes if there isn't a tmp directory to write error logs to
COPY --from=go-dev --chown=golang:root /tmp /tmp
COPY --from=go-dev --chown=golang:root /app/pinman /app/migrate /app/bin/
COPY --from=node-dev --chown=golang:root /usr/src/app/build /app/html

USER golang:root
EXPOSE 8080

ENV SPA_PATH=/app/html
WORKDIR "/app/bin"
ENTRYPOINT ["pinman"]
