# Variables
PROJECT_NAME := fullplex
CLIENT_DIR := ./client
SERVER_DIR := ./server
DOCKER_COMPOSE_FILE := ./docker-compose.yml

# Build the client binary
build-client:
	@echo "Building client binary..."
	cd $(CLIENT_DIR) && go build -o client client.go

# Build the server binary
build-server:
	@echo "Building server binary..."
	cd $(SERVER_DIR) && go build -o server server.go

# Build client Docker image
docker-build-client:
	@echo "Building client Docker image..."
	docker build -t $(PROJECT_NAME)-client $(CLIENT_DIR)

# Build server Docker image
docker-build-server:
	@echo "Building server Docker image..."
	docker build -t $(PROJECT_NAME)-server $(SERVER_DIR)

# Build all Docker images
docker-build-all: docker-build-client docker-build-server

# Run Docker Compose
docker-up:
	@echo "Starting Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

# Stop Docker Compose
docker-down:
	@echo "Stopping Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Clean up
clean:
	@echo "Cleaning up build artifacts..."
	rm -f $(CLIENT_DIR)/client $(SERVER_DIR)/server
	docker system prune -f

# Full build and deploy
deploy: docker-build-all docker-up

# Full clean up
clean-all: docker-down clean
