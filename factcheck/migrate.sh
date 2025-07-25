#!/bin/bash

# Database migration script for factcheck API
# Runs schema.sql against the PostgreSQL database

set -e

CONTAINER_NAME="factcheck-postgres"
DB_NAME="factcheck"
DB_USER="postgres"
DB_PASSWORD="postgres"
SCHEMA_FILE="data/postgres/schema.sql"

echo "Running database migration (DROP + CREATE)..."

# Check if schema file exists
if [ ! -f "$SCHEMA_FILE" ]; then
    echo "Error: Schema file not found at $SCHEMA_FILE"
    exit 1
fi

echo "Dropping existing tables..."

# Drop tables in reverse dependency order (due to foreign keys)
docker exec -i ${CONTAINER_NAME} psql -U ${DB_USER} -d ${DB_NAME} << EOF
DROP TABLE IF EXISTS user_messages CASCADE;
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS topics CASCADE;
DROP TABLE IF EXISTS messages_v2 CASCADE;
DROP TABLE IF EXISTS messages_v2_groups CASCADE;
EOF

echo "Creating fresh schema from $SCHEMA_FILE..."

# Execute the schema file
docker exec -i ${CONTAINER_NAME} psql -U ${DB_USER} -d ${DB_NAME} < ${SCHEMA_FILE}

echo "Migration completed successfully!"
echo ""
echo "Database tables created:"
echo "  - topics"
echo "  - messages" 
echo "  - user_messages"
echo ""
echo "You can now run your factcheck API with:"
echo "  cd factcheck && go run cmd/api/main.go" 