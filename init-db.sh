#!/bin/bash
set -e

# Check if database exists and create if not
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    SELECT 'CREATE DATABASE mockva'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'mockva')\gexec
EOSQL

echo "Database initialization check completed."
