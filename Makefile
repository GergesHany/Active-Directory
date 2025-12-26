.PHONY: help run up down reset-password build clean logs restart

# Default target
.DEFAULT_GOAL := help

# Variables
DOCKER_COMPOSE := docker compose
CONTAINER_NAME := samba-ad-dc
GO_BINARY := ad-example
PASSWORD := P@ssw0rd123!

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: up reset-password ## Build docker compose, reset password, and run the Go application
	@echo "Starting the application..."
	@sleep 3
	@$(MAKE) setup-test-data
	@go run main.go

up: ## Start the Samba AD Docker container
	@echo "Starting Samba AD container..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Waiting for Samba AD to initialize..."
	@sleep 8

down: ## Stop the Samba AD Docker container
	@echo "Stopping Samba AD container..."
	@$(DOCKER_COMPOSE) down

restart: down up reset-password ## Restart the Samba AD container and reset password
	@echo "Container restarted successfully"

reset-password: ## Reset the Administrator password in Samba AD
	@echo "Resetting Administrator password..."
	@docker exec $(CONTAINER_NAME) samba-tool user setpassword Administrator --newpassword='$(PASSWORD)' 2>/dev/null || \
		(echo "Waiting for container to be ready..." && sleep 5 && \
		docker exec $(CONTAINER_NAME) samba-tool user setpassword Administrator --newpassword='$(PASSWORD)')
	@echo "Password reset successfully"

build: ## Build the Go binary
	@echo "Building Go application..."
	@go build -o $(GO_BINARY) .
	@echo "Binary created: $(GO_BINARY)"

clean: ## Clean up Docker containers, volumes, and Go binary
	@echo "Cleaning up..."
	@$(DOCKER_COMPOSE) down -v
	@rm -f $(GO_BINARY)
	@echo "Cleanup complete"

logs: ## Show Samba AD container logs
	@docker logs $(CONTAINER_NAME) --tail 50 -f

status: ## Check the status of the Samba AD container
	@docker ps -a | grep $(CONTAINER_NAME) || echo "Container not found"

test-connection: ## Test LDAP connection without running the full app
	@echo "Testing LDAP connection..."
	@ldapsearch -x -H ldap://127.0.0.1:389 -b "DC=example,DC=com" -D "Administrator@example.com" -w "$(PASSWORD)" "(objectClass=*)" dn -LLL | head -10 || \
		echo "Connection test failed. Make sure the container is running and password is set."

setup-test-data: ## Create test users and groups in Active Directory
	@echo "Creating test user '$(shell grep TEST_USERNAME .env | cut -d'=' -f2)'..."
	@docker exec $(CONTAINER_NAME) samba-tool user create $$(grep TEST_USERNAME .env | cut -d'=' -f2) $$(grep TEST_PASSWORD .env | cut -d'=' -f2) --given-name=$$(grep TEST_FIRSTNAME .env | cut -d'=' -f2) --surname=$$(grep TEST_LASTNAME .env | cut -d'=' -f2) --mail-address=$$(grep TEST_EMAIL .env | cut -d'=' -f2) 2>/dev/null || echo "User might already exist"
	@echo "Creating test group '$$(grep TEST_GROUP .env | cut -d'=' -f2)'..."
	@docker exec $(CONTAINER_NAME) samba-tool group add $$(grep TEST_GROUP .env | cut -d'=' -f2) 2>/dev/null || echo "Group might already exist"
	@echo "Setting user attributes..."
	@docker exec $(CONTAINER_NAME) samba-tool user setexpiry $$(grep TEST_USERNAME .env | cut -d'=' -f2) --noexpiry 2>/dev/null || true
	@echo "Test data setup complete!"
