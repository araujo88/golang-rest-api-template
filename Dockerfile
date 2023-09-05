FROM golang:1.21.0-bookworm
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . /app
RUN swag init --dir cmd/server/
EXPOSE 8001
RUN CGO_ENABLED=1 go build -o bin/server cmd/server/main.go
CMD ./bin/server