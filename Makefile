setup:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./cmd/server/main.go -o ./docs
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

build:
	docker compose build --no-cache

run-local:
	docker start dockerPostgres
	docker start dockerRedis
	docker start dockerMongo
	export REDIS_HOST=localhost
	export POSTGRES_DB=go_app_dev
	export POSTGRES_USER=docker
	export POSTGRES_PASSWORD=password
	export POSTGRES_PORT=5435
	export JWT_SECRET_KEY=ObL89O3nOSSEj6tbdHako0cXtPErzBUfq8l8o/3KD9g=INSECURE
	export API_SECRET_KEY=cJGZ8L1sDcPezjOy1zacPJZxzZxrPObm2Ggs1U0V+fE=INSECURE
	export POSTGRES_HOST=localhost
	go run cmd/server/main.go

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
