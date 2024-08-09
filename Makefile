setup:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./cmd/server/main.go -o ./docs
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

build:
	docker compose build --no-cache

up:
	docker compose up

down:
	docker compose down

restart:
	docker compose restart

clean:
	docker stop go-rest-api-template
	docker stop dockerPostgres
	docker rm go-rest-api-template
	docker rm dockerPostgres
	docker rm dockerRedis
	docker image rm golang-rest-api-template-backend
	rm -rf .dbdata
