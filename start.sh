#!/bin/sh

set -e

echo "run db migration"
source /app/app.env # to read env in production
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"   # takes all param and exec
