#!/bin/bash
set -euo pipefail

COMPOSE_FILE="deploy/docker-compose.yml"
SERVICES=("$@")

if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker and try again."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

if [ ${#SERVICES[@]} -eq 0 ]; then
    echo "No services specified."
    exit 1
fi

echo "Building and starting services: ${SERVICES[*]}..."
docker-compose -f "$COMPOSE_FILE" up -d --build "${SERVICES[@]}"

echo "Migrating database..."
bash ./scripts/migrate.sh up

echo "Services started successfully"
