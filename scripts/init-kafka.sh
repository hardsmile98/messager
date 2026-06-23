#!/bin/bash
set -euo pipefail

KAFKA=(docker compose -f deploy/docker-compose.yml exec -it kafka)

TOPICS=(
    "messages:1:1"
)

for topic_config in "${TOPICS[@]}"; do
    IFS=':' read -r name partitions replication_factor <<< "$topic_config"
    echo "Creating topic: $name"
    "${KAFKA[@]}" opt/kafka/bin/kafka-topics.sh --create --if-not-exists \
        --bootstrap-server localhost:9092 \
        --topic "$name" \
        --partitions "$partitions" \
        --replication-factor "$replication_factor"
done
echo "Kafka topics created"
