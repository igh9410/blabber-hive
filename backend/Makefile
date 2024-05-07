include .env
export

# Makefile`

# Run the application
run:
	@echo "Running the application..."
	@go run cmd/blabber-hive/main.go

# Build the application using Docker
docker-push:
#@echo "Generating Swagger documentation before building the application..."
#	@swag init -g cmd/blabber-hive/main.go
#	@echo "Building the application using Docker..."
	@./scripts/docker-push.sh


# Run the application using Docker Compose
docker-run:
	@echo "Running the application using Docker Compose..."
	@./scripts/docker-run.sh

# Run golangci-lint
linter:
	@golangci-lint run

# Create database migration file
create-migration:
	@cd ./db/migrations && goose create $(filter-out $@,$(MAKECMDGOALS)) sql
	@echo "Migration created."

goose-version:
	@goose -dir db/migrations postgres "host=localhost user=$$POSTGRES_USERNAME password=$$POSTGRES_PASSWORD dbname=postgres sslmode=disable port=5432" version

# Run the database migrations
migrate-up:
	@echo "Running the database migrations..."
	@goose -dir db/migrations postgres "host=localhost user=$$POSTGRES_USERNAME password=$$POSTGRES_PASSWORD dbname=postgres sslmode=disable port=5432" up

# Rollback the database migrations
migrate-down:
	@echo "Rolling back the database migrations..."
	@goose -dir db/migrations postgres "host=localhost user=$$POSTGRES_USERNAME password=$$POSTGRES_PASSWORD dbname=postgres sslmode=disable port=5432" down

# Run the tests
test:
	@echo "Running the tests..."
	@go test ./... -v -cover -coverprofile=coverage.out
