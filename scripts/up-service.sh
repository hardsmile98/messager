#!/bin/bash
set -euo pipefail

COMPOSE_FILE="deploy/docker-compose.yml"
SERVICE="${1:-}"

if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker and try again."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

if [ -z "$SERVICE" ]; then
    echo "Usage: make up-service <name>"
    echo "Available services: auth chat message media"
    exit 1
fi

echo "Building and starting service: $SERVICE..."
docker-compose -f "$COMPOSE_FILE" up -d --build "$SERVICE"

echo "Service $SERVICE started successfully"
