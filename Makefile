FRONTEND_DIR = './web/app'

setup: mod yarn env-local

env: env-local keys-dev

env-local:
	@cp .env .env.local

yarn:
	@cd $(FRONTEND_DIR) &&	yarn

mod:
	@go mod download

clean:
	@rm -rf $(FRONTEND_DIR)/build

generate:
	@docker compose up --remove-orphans --build -d openapi-server openapi-client
	@rm -rf internal/app/generated web/app/src/api/generated && true
	@docker compose cp openapi-server:/out internal/app/generated
	@docker compose cp openapi-client:/out web/app/src/api/generated
	@docker compose stop openapi-server openapi-client

build: build-backend build-frontend

build-backend:
	@mkdir -p ./build
	@go build -v -o ./build/pinman ./cmd/pinman/

build-frontend:
	@cd $(FRONTEND_DIR) && yarn build

run: clean generate build-frontend run-migrate
	@SPA_PATH=./web/app/build go run cmd/pinman/main.go

run-backend: run-migrate
	@go run cmd/pinman/main.go

run-database:
	@docker compose up -d postgres

run-migrate: run-database
	@go run cmd/migrate/main.go

run-frontend:
	@cd $(FRONTEND_DIR) && REACT_APP_API_HOST=http://localhost:8080 yarn start

test: test-backend test-frontend

test-setup:
	@go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo

test-frontend:
	@cd $(FRONTEND_DIR) && yarn test

test-backend: test-setup
	@ginkgo	./...

keys-dev:
	@echo "Generating key-pair for jwt tokens..."
	@ssh-keygen -f `pwd`/token -t rsa -N '' -q
	@sed -i '' -e "s/TOKEN_PRIVATE_KEY.*/TOKEN_PRIVATE_KEY=`cat ./token | base64`/" "./.env.local"
	@sed -i '' -e "s/TOKEN_PUBLIC_KEY.*/TOKEN_PUBLIC_KEY=`cat ./token.pub | base64`/" "./.env.local"
	@rm -r ./token ./token.pub

test-backend-cov: test-setup
	@ginkgo --cover --race --json-report=report.json --output-dir=reports/go ./...