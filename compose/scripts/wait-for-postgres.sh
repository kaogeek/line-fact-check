#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
echo "wait-for-postgres.sh: HOST=$host"
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q'; do
  >&2 echo "wait-for-postgres.sh: PostgreSQL is unavailable - sleeping"
  sleep 1
done

>&2 echo "wait-for-postgres.sh: PostgreSQL is up!"
