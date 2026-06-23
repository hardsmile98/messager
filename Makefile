SERVICES := auth chat message media
REGISTRY := localhost:5000

up:
	bash ./scripts/up.sh

down:
	docker-compose -f deploy/docker-compose.yml down

migrate-up:
	bash ./scripts/migrate.sh up

migrate-down:
	bash ./scripts/migrate.sh down

init-kafka:
	bash ./scripts/init-kafka.sh