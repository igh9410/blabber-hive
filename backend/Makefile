include .env
PROTO_PATH="./proto"

# Makefile`
.PHONY: all run docker-push docker-run linter create-migration goose-version migrate-up migrate-down test

all: generate-docs generate-server 
# Run the application
run:
	@echo "Running the application..."
	@swag init -g cmd/blabber-hive/main.go
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
	@goose -dir internal/pkg/db/migrations create $(filter-out $@,$(MAKECMDGOALS)) sql
	@echo "Migration created."

goose-version:
	@goose -dir db/migrations postgres "host=localhost user=$$POSTGRES_USERNAME password=$$POSTGRES_PASSWORD dbname=postgres sslmode=disable port=5432" version

# Run the database migrations
migrate-up:
	@echo "Running the database migrations..."
	@goose -dir internal/pkg/db/migrations postgres "host=localhost user=$$POSTGRES_USERNAME password=$$POSTGRES_PASSWORD dbname=postgres sslmode=disable port=5432" up

# Rollback the database migrations
migrate-down:
	@echo "Rolling back the database migrations..."
	@goose -dir internal/pkg/db/migrations postgres "host=localhost user=$$POSTGRES_USERNAME password=$$POSTGRES_PASSWORD dbname=postgres sslmode=disable port=5432" down

# Run the tests
test:
	@echo "Running the tests..."
	@go test ./... -v -cover -coverprofile=coverage.out

# Generate Swagger documentation from proto files
generate-docs:
	@echo "Generating OpenAPI documentation..."
	@protoc --proto_path=${PROTO_PATH} -I=proto/third_party --openapi_out=api --openapi_opt=enum_type=string ${PROTO_PATH}/*.proto
# Generate server code from OpenAPI specification
generate-server:
	@echo "Generating server code from OpenAPI specification..."
	@oapi-codegen --generate=gin-server,strict-server,embedded-spec --package=api -o internal/api/server.gen.go api/openapi.yaml
	@oapi-codegen --generate=models --package=api -o internal/api/types.gen.go  api/openapi.yaml

# Generate client code from OpenAPI specification
generate-client:
	@echo "Generating client code from OpenAPI specification..."
	@oapi-codegen --generate=client -o internal/client/client.gen.go --package=client api/openapi.yaml



	