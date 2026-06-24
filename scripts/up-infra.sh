#!/bin/bash
set -euo pipefail

COMPOSE_FILE="deploy/docker-compose.yml"

if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker and try again."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

echo "Starting infrastructure..."
docker-compose -f "$COMPOSE_FILE" up -d postgres redis kafka

echo "Waiting for Kafka..."
until docker compose -f "$COMPOSE_FILE" exec -T kafka opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server localhost:9092 &>/dev/null; do
    sleep 1
done

echo "Creating Kafka topics..."
bash ./scripts/init-kafka.sh

echo "Infrastructure started successfully"
