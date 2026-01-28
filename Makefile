.PHONY: help build run test test-unit test-integration clean deps keploy-record keploy-test keploy-list

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/todo-server main.go

run: ## Run the application
	go run main.go

test: ## Run all tests (unit + integration)
	go test ./... -v

test-unit: ## Run only unit tests
	go test ./handlers/... ./repository/... -v

test-integration: ## Run integration tests (requires database)
	go test -v -run TestIntegration

clean: ## Clean build artifacts and test cache
	rm -rf bin/
	go clean -testcache

keploy-record: ## Record mocks with Keploy
	keploy record -c "go run main.go"

keploy-test: ## Run tests with Keploy mocks
	keploy test -c "go test ./... -v"

keploy-list: ## List available Keploy mock sets
	@echo "Available mock sets:"
	@ls -la keploy/ 2>/dev/null || echo "No mocks found. Run 'make keploy-record' first."
