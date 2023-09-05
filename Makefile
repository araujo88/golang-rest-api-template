setup:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --dir cmd/server/
	go build -o bin/server cmd/server/main.go

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
