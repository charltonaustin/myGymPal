#!/bin/bash
# Creates additional PostgreSQL databases listed in POSTGRES_MULTIPLE_DATABASES.
# Usage: set POSTGRES_MULTIPLE_DATABASES=db1,db2 in the container environment.
set -e

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
  for db in $(echo "$POSTGRES_MULTIPLE_DATABASES" | tr ',' ' '); do
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
      CREATE DATABASE $db;
      GRANT ALL PRIVILEGES ON DATABASE $db TO $POSTGRES_USER;
EOSQL
  done
fi
