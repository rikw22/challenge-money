.PHONY: run test test-coverage build

run:
	@echo "Starting postgres container..."
	@docker compose up -d
	@echo "Container started!\n"
	@if command -v air > /dev/null 2>&1; then \
		echo "Running with air (live reloading enabled)..."; \
		air; \
	else \
		echo "air command is not available. Falling back to go run without live reloading..."; \
		echo "To enable live reloading, install air: go install github.com/air-verse/air@latest"; \
		go run ./cmd/api/main.go; \
	fi

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build:
	go build

docker-up:
	docker compose up -d

docker-down:
	docker compose down