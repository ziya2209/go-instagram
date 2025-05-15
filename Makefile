.PHONY: build up down restart logs ps clean redeploy help

# Default command, shows help
.DEFAULT_GOAL := help

# Build or rebuild services
build:
	docker-compose build

# Create and start containers
up:
	docker-compose up -d

# Stop and remove containers, networks
down:
	docker-compose down

# Restart services
restart:
	docker-compose restart

# View output from containers
logs:
	docker-compose logs -f

# List containers
ps:
	docker-compose ps

# Remove stopped containers and unused images
clean:
	docker-compose down --rmi local
	docker system prune -f

# Redeploy changes (rebuild and restart)
redeploy:
	docker-compose down
	docker-compose build
	docker-compose up -d

# Redeploy only the app service (faster for code changes)
ra:
	docker-compose stop app
	docker-compose rm -f app
	docker-compose build app
	docker-compose up -d app

# Show help message
help:
	@echo "Instagram Clone Development Commands:"
	@echo ""
	@echo "make build        - Build or rebuild services"
	@echo "make up           - Create and start containers in detached mode"
	@echo "make down         - Stop and remove containers and networks"
	@echo "make restart      - Restart all services"
	@echo "make logs         - View output from containers"
	@echo "make ps           - List running containers"
	@echo "make clean        - Remove stopped containers and unused images"
	@echo "make redeploy     - Perform a full rebuild and restart of all services"
	@echo "make redeploy-app - Rebuild and restart only the app service"
	@echo "make help         - Show this help message"
	@echo ""
