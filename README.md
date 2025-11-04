# grpc_service

> gRPC service implemented in Go

## Overview
This repository contains a gRPC service implemented in Go.  
The project includes:  
- Server implementation (`server/`)  
- Client example (`client/`)  
- API definitions in Protobuf (`proto/`)  
- Go module configuration (`go.mod`, `go.sum`)
- Docker setup for PostgreSQL database

## Repository Structure

├── proto/  .proto files defining services and messages \
├── server/  gRPC server implementation \
├── client/  Example client using the service \
├── go.mod   \
├── go.sum   \
└── .gitignore

## Requirements
- Go 1.20 or higher  
- Protocol Buffers compiler (`protoc`)  
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins for generating Go code from `.proto` files
- Docker (for running PostgreSQL)

## Docker Setup
The service requires a PostgreSQL database. You can run it using Docker:
docker run -d --name go_grps_postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres123 -e POSTGRES_DB=go_grps -p 5434:5432 postgres
This command will:

Start a PostgreSQL container named go_grps_postgres
Set the username to postgres and password to postgres123
Create a database named go_grps
Map port 5434 on your host to 5432 in the container

## Installation and Running
#### 1. Running the Services

Run the server
cd server
go run main.go

Run the client
cd client
go run main.go

The client will connect to the server and perform an example gRPC request.

Generating Code from Protobuf
If you modify or add .proto files, generate the Go code with:
protoc --go_out=. --go-grpc_out=. proto/book.proto
Make sure protoc-gen-go and protoc-gen-go-grpc are in your $PATH.
