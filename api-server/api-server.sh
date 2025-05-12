#!/bin/bash

CONTAINER_NAME="pg-container"
POSTGRES_PASSWORD="secret"
DB_NAME="gopgtest"
APP_PATH="cmd/api-server/main.go"

# Pull postgres image
echo "🔄 Pulling postgres image (if not already present)..."
docker pull postgres

# Check if container already exists
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    # If it's not running, start it
    if [ "$(docker ps -q -f status=running -f name=$CONTAINER_NAME)" ]; then
        echo "✅ Container $CONTAINER_NAME is already running."
    else
        echo "▶️ Starting existing container $CONTAINER_NAME..."
        docker start $CONTAINER_NAME
    fi
else
    echo "🚀 Creating and starting new postgres container..."
    docker run --name $CONTAINER_NAME -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD -p 5432:5432 -d postgres
fi

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to become ready..."
until docker exec -u postgres $CONTAINER_NAME pg_isready -q; do
    sleep 1
done

# Create the database if it doesn't exist
echo "🛠️ Creating database '$DB_NAME' if it doesn't exist..."
docker exec -u postgres $CONTAINER_NAME psql -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
    docker exec -u postgres $CONTAINER_NAME createdb -U postgres $DB_NAME

echo "✅ Database '$DB_NAME' is ready."

# Run your Go application
echo "🚀 Starting Go application at '$APP_PATH'..."
go run "$APP_PATH"
