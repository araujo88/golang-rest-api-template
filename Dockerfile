FROM golang:1.21.0-bookworm
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod vendor
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/server/main.go -o ./docs
RUN CGO_ENABLED=1 go build -o bin/server cmd/server/main.go
CMD ./bin/server