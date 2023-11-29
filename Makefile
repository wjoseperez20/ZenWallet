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