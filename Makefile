setup:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init

build:
	docker compose up --build

up:
	docker compose up

down:
	docker compose down

restart:
	docker compose restart

rebuild:
	docker compose up --build

clean:
	docker stop go-rest-api-template
	docker stop dockerPostgres
	docker rm go-rest-api-template
	docker rm dockerPostgres
	docker image rm go-rest-api-template
	rm -rf .dbdata