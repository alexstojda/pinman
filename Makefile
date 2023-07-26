FRONTEND_DIR = './web/app'

setup: setup-frontend setup-backend

setup-frontend:
	@cd $(FRONTEND_DIR) && yarn

setup-backend: mod

mod:
	@go mod download

env: env-local keys

env-local:
	@cp .env .env.local

clean:
	@rm -rf $(FRONTEND_DIR)/build

generate: generate-backend generate-frontend

generate-frontend:
	@docker compose up --remove-orphans --build -d openapi-client
	@rm -rf web/app/src/api/generated && true
	@docker compose cp openapi-client:/out web/app/src/api/generated
	@docker compose stop openapi-client

generate-backend:
	@docker compose up --remove-orphans --build -d openapi-server
	@rm -rf internal/app/generated && true
	@docker compose cp openapi-server:/out internal/app/generated
	@docker compose stop openapi-server

build: build-backend build-frontend

build-backend:
	@mkdir -p ./build
	@go build -v -o ./build/pinman main.go

build-frontend:
	@cd $(FRONTEND_DIR) && yarn build

run: clean generate build-frontend
	@SPA_PATH=./web/app/build go run main.go

run-backend:
	@go run main.go

run-database:
	@docker compose up -d postgres

run-frontend:
	@cd $(FRONTEND_DIR) && REACT_APP_API_HOST=http://localhost:8080 yarn start

test: test-backend test-frontend

test-setup:
	@go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo
	@go install -mod=mod github.com/vektra/mockery/v2@v2.32.0
	@mockery

test-frontend:
	@cd $(FRONTEND_DIR) && yarn test
	@rm -rf ./reports/ts && true
	@mkdir -p ./reports/ts
	@mv -f $(FRONTEND_DIR)/coverage $(FRONTEND_DIR)/reports ./reports/ts/

test-backend: test-setup
	@ginkgo	./...

keys:
	@echo "Generating key-pair for jwt tokens..."
	@ssh-keygen -f `pwd`/token -t rsa -N '' -q
	@sed -i "s/TOKEN_PRIVATE_KEY.*/TOKEN_PRIVATE_KEY=`cat ./token | base64 | tr -d '\n'`/" "./.env.local"
	@sed -i "s/TOKEN_PUBLIC_KEY.*/TOKEN_PUBLIC_KEY=`cat ./token.pub | base64 | tr -d '\n'`/" "./.env.local"
	@rm -r ./token ./token.pub

test-backend-cov: test-setup
	@ginkgo --cover \
		--race \
		--json-report=report.json \
		--output-dir=reports/go \
		--skip-package generated \
		--skip-file mock \
		./...
