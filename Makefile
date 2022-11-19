FRONTEND_DIR = './web/app'

setup: mod yarn env-local

env-local: keys-dev
	@cp .env .env.local

yarn:
	@cd $(FRONTEND_DIR) &&	yarn

mod:
	@go mod download

clean:
	@rm -rf $(FRONTEND_DIR)/build

build: build-backend build-frontend

build-backend:
	@mkdir -p ./build
	@go build -v -o ./build/pinman ./cmd/pinman/

build-frontend:
	@cd $(FRONTEND_DIR) && yarn build

run: clean build-frontend run-migrate
	@go run cmd/pinman/main.go

run-backend:
	@go run cmd/pinman/main.go

run-database:
	@docker compose up -d postgres

run-migrate: run-database
	@go run cmd/migrate/main.go

run-frontend:
	@cd $(FRONTEND_DIR) && REACT_APP_API_HOST=http://localhost:8080 yarn start

test:
	@go test ./...

keys-dev:
	@echo "Generating key-pair for access_token..."
	@ssh-keygen -f `pwd`/access_token -t rsa -N '' -q
	@sed -i '' -e "s/ACCESS_TOKEN_PRIVATE_KEY.*/ACCESS_TOKEN_PRIVATE_KEY=`cat ./access_token | base64`/" "./.env.local"
	@sed -i '' -e "s/ACCESS_TOKEN_PUBLIC_KEY.*/ACCESS_TOKEN_PUBLIC_KEY=`cat ./access_token.pub | base64`/" "./.env.local"
	@rm -r ./access_token ./access_token.pub
	@echo "Generating key-pair for refresh_token..."
	@ssh-keygen -f `pwd`/refresh_token -t rsa -N '' -q
	@sed -i '' -e "s/REFRESH_TOKEN_PRIVATE_KEY.*/REFRESH_TOKEN_PRIVATE_KEY=`cat ./refresh_token | base64`/" "./.env.local"
	@sed -i '' -e "s/REFRESH_TOKEN_PUBLIC_KEY.*/REFRESH_TOKEN_PUBLIC_KEY=`cat ./refresh_token.pub | base64`/" "./.env.local"
	@rm -r ./refresh_token ./refresh_token.pub

#test-cov:
#	mkdir -p coverage
#	@go test -covermode=atomic -coverprofile=./coverage/coverage.txt ./...
#	@go get github.com/axw/gocov/gocov
#	@go get github.com/AlekSi/gocov-xml
#	@gocov convert ./coverage/coverage.txt | gocov-xml > ./coverage/coverage.xml