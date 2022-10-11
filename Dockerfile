FROM golang:1.19.2-bullseye
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . /app
EXPOSE 8001
CMD go run main.go

