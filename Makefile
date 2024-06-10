help:
	@echo "Use 'make <target>'"
	@echo "  dev        		to run docker compose with development environment"
	@echo "  prod       		to run docker compose without development environment"
	@echo "  generate-docs  	to generate swagger (OpenAPI 2.0) docs"
	@echo "  test       		to run tests"
	@echo "  e2e-test   		to run end-to-end tests"
	@echo "  build-prod       	to build the production image"
	@echo "  push-prod       	to push the production image"


# Run docker compose with development environment
dev:
	docker compose -f docker-compose.dev.yml up --build

# Run docker compose without development environment
prod:
	docker build . --tag bisquitt_psk && \
	docker compose up --build

# Generate swagger (OpenAPI 2.0) docs
generate-docs:
	docker compose -f docker-compose.dev.yml run --rm api go install github.com/swaggo/swag/cmd/swag@latest && swag init -d "./cmd/api","./pkg"

# Run tests
test:
	go test -race -v ./pkg/...

# Run end-to-end tests
e2e-test:
	@TEST=1 docker build . --tag bisquitt_psk && \
    (timeout 60 docker-compose -f ./tests/docker-compose.yml up --exit-code-from test_runner || true) && \
    (docker-compose -f ./tests/docker-compose.yml down; exit $$?)

# Build the production image with latest tag and version tag and push it to the registry
build-push-prod:
	docker build . --tag ghcr.io/energostack/bisquitt_psk:latest && \
	docker build . --tag ghcr.io/energostack/bisquitt_psk:$(shell git describe --tags --abbrev=0)
	docker push ghcr.io/energostack/bisquitt_psk:latest && \
	docker push ghcr.io/energostack/bisquitt_psk:$(shell git describe --tags --abbrev=0)










