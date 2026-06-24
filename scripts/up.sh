#!/bin/bash
set -euo pipefail

bash ./scripts/up-infra.sh
bash ./scripts/up-services.sh "$@"

echo "All services started successfully"
