FROM golang:1.21.0-bookworm
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY . /app
RUN go mod download
RUN go mod vendor
EXPOSE 8001
RUN CGO_ENABLED=1 go build -o bin/server cmd/server/main.go
CMD ./bin/server