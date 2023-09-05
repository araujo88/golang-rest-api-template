# golang-rest-api-template

## Overview

This repository provides a template for building a RESTful API using Go with features like JWT Authentication, rate limiting, Swagger documentation, and database operations using GORM. The application uses the Gin Gonic web framework and is containerized using Docker.

## Features

- RESTful API endpoints for CRUD operations.
- JWT Authentication.
- Rate Limiting.
- Swagger Documentation.
- PostgreSQL database integration using GORM.
- Redis cache.
- Dockerized application for easy setup and deployment.

## Folder structure

```
golang-rest-api-template/
|-- bin/
|-- cmd/
|   |-- server/
|       |-- main.go
|-- pkg/
|   |-- api/
|       |-- handler.go
|       |-- router.go
|   |-- models/
|       |-- user.go
|   |-- database/
|       |-- db.go
|-- scripts/
|-- Dockerfile
|-- go.mod
|-- go.sum
|-- README.md
```

### Explanation of Directories and Files:

1. **`bin/`**: Contains the compiled binaries.

2. **`cmd/`**: Main applications for this project. The directory name for each application should match the name of the executable.

    - **`main.go`**: The entry point.

3. **`pkg/`**: Libraries and packages that are okay to be used by applications from other projects. 

    - **`api/`**: API logic.
        - **`handler.go`**: HTTP handlers.
        - **`router.go`**: Routes.
    - **`models/`**: Data models.
    - **`database/`**: Database connection and queries.

4. **`scripts/`**: Various build, install, analysis, etc., scripts.

## Getting Started

### Prerequisites

- Go 1.15+
- Docker
- Docker Compose

### Installation

1. Clone the repository

```bash
git clone https://github.com/araujo88/golang-rest-api-template
```

2. Navigate to the directory

```bash
cd golang-rest-api-template
```

3. Build and run the Docker containers

```bash
make setup && make build && make up
```

### Environment Variables

You can set the environment variables in the `.env` file. Here are some important variables:

- `POSTGRES_HOST`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `JWT_SECRET`
- `API_SECRET_KEY`

### API Documentation

The API is documented using Swagger and can be accessed at:

```
http://localhost:8001/swagger/index.html
```

## Usage

### Endpoints

- `GET /api/v1/books`: Get all books.
- `GET /api/v1/books/:id`: Get a single book by ID.
- `POST /api/v1/books`: Create a new book.
- `PUT /api/v1/books/:id`: Update a book.
- `DELETE /api/v1/books/:id`: Delete a book.
- `POST /api/v1/login`: Login.
- `POST /api/v1/register`: Register a new user.

### Authentication

To use authenticated routes, you must include the `Authorization` header with the JWT token.

```bash
curl -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8001/api/v1/books
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## TODOs

 - Unit tests
 - e2e tests
