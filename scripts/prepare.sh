#!/usr/bin/env bash
set -e

echo "🔧 Preparing environment..."

DB_HOST="localhost"
DB_PORT="5432"
DB_USER="validator"
DB_PASS="val1dat0r"
DB_NAME="project-sem-1"

export PGPASSWORD="$DB_PASS"

echo "Waiting for PostgreSQL..."
for i in {1..30}; do
  if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "select 1" >/dev/null 2>&1; then
    echo "✅ PostgreSQL is ready"
    break
  fi
  echo "waiting..."
  sleep 1
done


echo "Database prepared successfully"
