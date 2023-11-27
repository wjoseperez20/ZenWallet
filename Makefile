<<<<<<< d9f4b386e3f7cc54698cc25094461bf7f54876f2
PROJECT_NAME := "zenwallet"
DOCKER_IMAGE_VERSION := $(if $(CI_COMMIT_TAG),$(CI_COMMIT_TAG),$(CI_COMMIT_SHORT_SHA))

.PHONY: help all fmt tidy lint test coverage run build clean build-docker-local docker-run-local docker-cleanup

help: ## Displays this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: fmt tidy lint test

fmt: ## Runs go fmt
	go fmt ./...

tidy: ## Runs go mod tidy to update dependencies
	go mod tidy

lint: ## Runs linter
	./bin/lint.sh;

test: ## Runs unit tests
	./bin/test.sh;

release: ## Update CHANGELOG.md file
	./bin/release.sh

coverage: ## Generates global code coverage report
	./bin/coverage.sh;

run: ## Runs the application
	@go run cmd/server/main.go

build: clean ## Builds the binary file
	mkdir -p build && go build -v cmd/server/main.go > build/${PROJECT_NAME}

clean: ## Removes previous build
	@rm -rf build

docker-build-local: ## Build docker images locally (via docker-compose)
	@docker-compose build

docker-run-local: ## Runs the app locally using docker-compose. Server will be listening on http://localhost:8080
	@docker-compose up --build

docker-cleanup: ## Removes all docker containers created by docker-run-local
	@docker-compose rm --stop
=======
PROJECT_NAME := "Zenwallet"

#Build Documentation
setup:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	swag i -g cmd/server/main.go -o docs
	go build -o bin/server cmd/server/main.go

# Docker compose build will build the images if they don't exist.
build:
	docker compose build --no-cache

#Docker compose up will start the containers in the background and leave them running.
up:
	docker compose up

database:
	 bash ./scripts/database_structure.sh

#Docker compose down will stop the containers and remove them.
down:
	docker compose down

#Docker compose restart will restart the containers.
restart:
	docker compose restart

#Clean will stop and remove the containers and images.
clean:
	docker stop Zenwallet
	docker stop Postgres
	docker stop Redis
	docker rm Zenwallet
	docker rm Postgres
	docker rm Redis
	docker image rm zenwallet-api
>>>>>>> feat: containers work
