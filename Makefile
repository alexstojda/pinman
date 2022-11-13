setup: mod
	yarn install

mod:
	go mod download

build: build-backend build-frontend

build-backend:
	mkdir -p ./build
	go build -v -o ./build/pinman ./cmd/pinman/

build-frontend:
	yarn build

run: build-frontend
	SPA_PATH=./build go run cmd/pinman/main.go

run-backend:
	go run cmd/pinman/main.go

run-frontend:
	REACT_APP_API_HOST=http://localhost:8080 yarn start

test:
	@go test ./...

#test-cov:
#	mkdir -p coverage
#	@go test -covermode=atomic -coverprofile=./coverage/coverage.txt ./...
#	@go get github.com/axw/gocov/gocov
#	@go get github.com/AlekSi/gocov-xml
#	@gocov convert ./coverage/coverage.txt | gocov-xml > ./coverage/coverage.xml