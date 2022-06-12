#!/bin/bash
SERVER_ADDRESS=localhost \
SERVER_PORT=8181 \
DB_USER=root \
DB_PASSWORD=password \
DB_HOST=localhost \
DB_PORT=3306 \
DB_NAME=banking \
JWT_SECRET=secret \
go run main.go