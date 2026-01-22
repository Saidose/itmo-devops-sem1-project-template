#!/usr/bin/env bash
set -e

export HTTP_PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=validator
export DB_PASS=val1dat0r
export DB_NAME=project-sem-1

echo "Starting server"
nohup go run ./cmd/api > server.log 2>&1 &
echo "waiting for healthy"
for i in $(seq 1 15); do
  if curl -fsS "http://localhost:8080/health" >/dev/null 2>&1; then
    echo "[run] OK"
    exit 0
  fi
  sleep 2
done

echo "ERROR: healthcheck timeout:" >&2
tail -n 150 "server.log" >&2 || true
exit 1