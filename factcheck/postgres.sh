#!/bin/bash

# PostgreSQL Docker container for factcheck API
# Configuration matches settings in cmd/api/config/config.go

set -e

CONTAINER_NAME="factcheck-postgres"
IMAGE_NAME="postgres:15"
HOST_PORT=5432
CONTAINER_PORT=5432
DB_NAME="factcheck"
DB_USER="postgres"
DB_PASSWORD="postgres"

echo "Starting PostgreSQL container for factcheck API..."

# Check if container already exists
if docker ps -a --format "table {{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
    echo "Container ${CONTAINER_NAME} already exists."
    
    # Check if it's running
    if docker ps --format "table {{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
        echo "Container is already running."
        echo "Connection details:"
        echo "  Host: localhost"
        echo "  Port: ${HOST_PORT}"
        echo "  Database: ${DB_NAME}"
        echo "  User: ${DB_USER}"
        echo "  Password: ${DB_PASSWORD}"
        exit 0
    else
        echo "Starting existing container..."
        docker start ${CONTAINER_NAME}
    fi
else
    echo "Creating new PostgreSQL container..."
    docker run -d \
        --name ${CONTAINER_NAME} \
        -e POSTGRES_DB=${DB_NAME} \
        -e POSTGRES_USER=${DB_USER} \
        -e POSTGRES_PASSWORD=${DB_PASSWORD} \
        -p ${HOST_PORT}:${CONTAINER_PORT} \
        -v factcheck_postgres_data:/var/lib/postgresql/data \
        ${IMAGE_NAME}
fi

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker exec ${CONTAINER_NAME} pg_isready -U ${DB_USER} -d ${DB_NAME}; do
    echo "PostgreSQL is not ready yet. Waiting..."
    sleep 2
done

echo "PostgreSQL is ready!"
echo ""
echo "Connection details:"
echo "  Host: localhost"
echo "  Port: ${HOST_PORT}"
echo "  Database: ${DB_NAME}"
echo "  User: ${DB_USER}"
echo "  Password: ${DB_PASSWORD}"
echo ""
echo "To connect with psql:"
echo "  psql -h localhost -p ${HOST_PORT} -U ${DB_USER} -d ${DB_NAME}"
echo ""
echo "To stop the container:"
echo "  docker stop ${CONTAINER_NAME}"
echo ""
echo "To remove the container:"
echo "  docker rm ${CONTAINER_NAME}" 