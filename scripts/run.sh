#!/usr/bin/env bash
set -e

export HTTP_PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=validator
export DB_PASS=val1dat0r
export DB_NAME=project-sem-1

echo "Starting server"

go run ./cmd/server
