#!/bin/bash
set -euo pipefail

if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker and try again."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

echo "Starting services..."
docker-compose -f deploy/docker-compose.yml up -d

echo "Waiting for services to start..."
until docker compose -f deploy/docker-compose.yml exec -it kafka opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server localhost:9092 &>/dev/null; do
    sleep 1
done

echo "Creating Kafka topics..."
bash ./scripts/init-kafka.sh

echo "Migrating database..."
bash ./scripts/migrate.sh up

echo "Services started successfully"