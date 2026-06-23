#!/bin/bash
set -euo pipefail

ACTION="${1:-up}"
MIGRATION_DIR="deploy/migrations"
DB_URL="postgres://postgres:postgres@localhost:5432/messager?sslmode=disable"
MIGRATE="migrate"

if [ "$ACTION" != "up" ] && [ "$ACTION" != "down" ]; then
  echo "Invalid action: $ACTION"
  echo "Usage: $0 [up|down]"
  exit 1
fi

if [ ! -d "$MIGRATION_DIR" ]; then
  echo "Migration directory not found: $MIGRATION_DIR"
  exit 1
fi

echo "Applying migrations ($ACTION)..."
"$MIGRATE" -database "$DB_URL" -path "$MIGRATION_DIR" "$ACTION"
echo "Migrations completed"