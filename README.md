# grpc_service

> gRPC service implemented in Go with PostgreSQL support

---

## 📘 Overview

This repository contains a **gRPC microservice** written in **Go**, designed to demonstrate a clean architecture with a PostgreSQL database.  

It includes:  
- Server implementation (`server/`)  
- Client example (`client/`)  
- API definitions in Protobuf (`proto/`)  
- Configuration via `go.mod`, `.env`, and `config.yml`  
- Docker setup for PostgreSQL database  

---

## Repository Structure

```bash
├── proto/              # .proto files defining services and messages
├── server/             
│   ├── grpc_server.go   
│   ├── cmd/            # main.go (entry point)
│   ├── configs/        # config.yml
│   ├── models/         # models.go 
│   └── pkg/            # handler, repository, service           
├── client/             # Example gRPC client
├── go.mod   
├── go.sum
├── .env   
└── .gitignore
````

---

## Requirements

* Go **1.20+**
* Protocol Buffers compiler (`protoc`)
* Go plugins:

  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```
* Docker (for PostgreSQL setup)

---

## Docker Setup

The service requires a PostgreSQL database.
Run the following command to create a Docker container:

```bash
docker run -d \
  --name go_grps_postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres123 \
  -e POSTGRES_DB=go_grps \
  -p 5434:5432 \
  postgres
```

This will:

* Start a PostgreSQL container named **go_grps_postgres**
* Create a database **go_grps**
* Use user: `postgres`, password: `postgres123`
* Map port **5434** on host → **5432** in container

---

## API Client Endpoints

### Authentication

| Method | Path            | Description         |
| ------ | --------------- | ------------------- |
| `POST` | `/auth/sign-up` | Register a new user |
| `POST` | `/auth/sign-in` | Authenticate a user |

### Books

| Method   | Path         | Description           |
| -------- | ------------ | --------------------- |
| `POST`   | `/books/`    | Create a new book     |
| `GET`    | `/books/`    | Retrieve all books    |
| `GET`    | `/books/:id` | Retrieve a book by ID |
| `PUT`    | `/books/:id` | Update a book by ID   |
| `DELETE` | `/books/:id` | Delete a book by ID   |

---

## Running the Services

### Run the Server

```bash
go run server/cmd/main.go
```

### Run the Client

```bash
go run client/main.go
```

The client will connect to the gRPC server and perform example requests.

---

## Generating Code from Protobuf

If you modify or add `.proto` files, generate the corresponding Go code using:

```bash
protoc --go_out=. --go-grpc_out=. proto/book.proto
```


