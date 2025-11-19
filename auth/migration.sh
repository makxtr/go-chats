#!/bin/bash
source ./local.env

export MIGRATION_DSN="host=${DB_HOST:-postgres-auth} port=5432 dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_DIR:-migrations}" postgres "${MIGRATION_DSN}" up -v
