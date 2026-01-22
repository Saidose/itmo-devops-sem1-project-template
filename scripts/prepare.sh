#!/usr/bin/env bash
set -e

echo "Preparing environment..."

if ! command -v psql >/dev/null 2>&1; then
  echo "PostgreSQL not installed"
  exit 1
fi

echo "PostgreSQL detected"

DB_USER="validator"
DB_PASS="val1dat0r"
DB_NAME="project-sem-1"
DB_PORT="5432"


sudo -u postgres psql <<EOF
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '${DB_USER}') THEN
      CREATE ROLE ${DB_USER} LOGIN PASSWORD '${DB_PASS}';
   END IF;
END
\$\$;
EOF

sudo -u postgres psql <<EOF
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '${DB_NAME}') THEN
      CREATE DATABASE "${DB_NAME}" OWNER ${DB_USER};
   END IF;
END
\$\$;
EOF


echo "Database prepared successfully"
