#!/bin/bash
source ./local.env

export MIGRATION_DSN="host=${DB_HOST:-postgres-auth} port=5432 dbname=$PG_DATABASE_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_DIR:-migrations}" postgres "${MIGRATION_DSN}" up -v
