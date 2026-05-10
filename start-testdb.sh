#!/bin/bash

# Script to start PostgreSQL test database with proper configuration

CONTAINER_NAME="testdb"
POSTGRES_PASSWORD="testpswd"
PORT="5433"
IMAGE="postgres"

echo "🚀 Starting PostgreSQL test database..."

# Check if container exists
if [ "$(docker ps -aq -f name=^${CONTAINER_NAME}$)" ]; then
    # Container exists, check if it's running
    if [ "$(docker ps -q -f name=^${CONTAINER_NAME}$)" ]; then
        echo "✓ Container '${CONTAINER_NAME}' is already running."
    else
        echo "Starting existing container '${CONTAINER_NAME}'..."
        docker start ${CONTAINER_NAME}
    fi
else
    # Container doesn't exist, create and run it
    echo "Creating new container '${CONTAINER_NAME}'..."
    docker run --name ${CONTAINER_NAME} \
        -p ${PORT}:5432 \
        -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
        -e POSTGRES_DB=postgres \
        -d ${IMAGE}

    # Wait for PostgreSQL to be ready
    echo "Waiting for PostgreSQL to be ready..."
    sleep 3

    # Wait for the database to accept connections
    for i in {1..30}; do
        if docker exec ${CONTAINER_NAME} pg_isready -U postgres > /dev/null 2>&1; then
            echo "✓ PostgreSQL is ready!"
            break
        fi
        if [ $i -eq 30 ]; then
            echo "✗ Timeout waiting for PostgreSQL to be ready"
            exit 1
        fi
        sleep 1
    done
fi

# Check if the container is running
if [ "$(docker ps -q -f name=^${CONTAINER_NAME}$)" ]; then
    echo "✓ Container '${CONTAINER_NAME}' is running on port ${PORT}"

    # Check if database setup is needed
    echo ""
    echo "📝 Setting up database schema and test data..."

    # Run the setup SQL script
    if [ -f "setup_postgres.sql" ]; then
        docker exec -i ${CONTAINER_NAME} psql -U postgres -d postgres < setup_postgres.sql

        if [ $? -eq 0 ]; then
            echo "✓ Database setup completed successfully!"
            echo ""
            echo "📊 Database Information:"
            echo "  Host: localhost"
            echo "  Port: ${PORT}"
            echo "  Database: testdb"
            echo "  User: testuser"
            echo "  Password: testpass"
            echo ""
            echo "  Admin User: postgres"
            echo "  Admin Password: ${POSTGRES_PASSWORD}"
            echo ""
            echo "  Connection string: postgresql://testuser:testpass@localhost:${PORT}/testdb"
            echo ""
            echo "🎯 You can now run: ./gql-validate check"
        else
            echo "⚠️  Database setup failed. You may need to run it manually:"
            echo "  docker exec -i ${CONTAINER_NAME} psql -U postgres -d postgres < setup_postgres.sql"
        fi
    else
        echo "⚠️  setup_postgres.sql not found. Database schema not created."
        echo "  Connection string: postgresql://postgres:${POSTGRES_PASSWORD}@localhost:${PORT}/postgres"
    fi
else
    echo "✗ Failed to start container '${CONTAINER_NAME}'"
    exit 1
fi
