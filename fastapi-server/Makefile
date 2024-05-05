# Makefile

# Run the application
run:
	@echo "Running the application..."
	@poetry run uvicorn app.main:app --reload

# Build the application using Docker
docker-build:
	@echo "Building the application using Docker..."
	@./scripts/docker-build.sh

# Run the application using Docker Compose
docker-run:
	@echo "Running the application using Docker Compose..."
	@./scripts/docker-run.sh

# Run flake8 linter
lint:
	@echo "Running flake8 linter..."
	@poetry run flake8 app/

# Run black code formatter
format:
	@echo "Running black code formatter..."
	@poetry run black app/

# Run the tests
test:
	@echo "Running the tests..."
	@poetry run pytest tests/ -v --cov=app --cov-report=term-missing

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	@poetry run coverage html

# Create a new migration file using Alembic
create-migration:
	@echo "Creating a new migration file..."
	@poetry run alembic revision --autogenerate -m "$(filter-out $@,$(MAKECMDGOALS))"

# Run the database migrations
migrate:
	@echo "Running the database migrations..."
	@poetry run alembic upgrade head

# Rollback the last database migration
migrate-rollback:
	@echo "Rolling back the last database migration..."
	@poetry run alembic downgrade -1

# Generate API documentation using Sphinx
docs:
	@echo "Generating API documentation..."
	@cd docs && poetry run make html
