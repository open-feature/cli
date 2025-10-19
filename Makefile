.PHONY: help
help:
	@echo "Available commands:"
	@echo "  help                     - Show this help message"
	@echo "  test                     - Run unit tests"
	@echo "  test-integration         - Run all integration tests"
	@echo "  test-integration-csharp  - Run C# integration tests"
	@echo "  test-integration-go      - Run Go integration tests"
	@echo "  test-integration-nodejs  - Run NodeJS integration tests"
	@echo "  generate                 - Generate all code (API clients, docs, schema)"
	@echo "  generate-api             - Generate API clients from OpenAPI specs"
	@echo "  generate-docs            - Generate documentation"
	@echo "  generate-schema          - Generate schema"
	@echo "  fmt                      - Format Go code"

.PHONY: test
test: 
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests passed successfully!"

# Dagger-based integration tests
.PHONY: test-integration-csharp
test-integration-csharp:
	@echo "Running C# integration test with Dagger..."
	@go run ./test/integration/cmd/csharp/run.go

.PHONY: test-integration-go
test-integration-go:
	@echo "Running Go integration test with Dagger..."
	@go run ./test/integration/cmd/go/run.go

.PHONY: test-integration-nodejs
test-integration-nodejs:
	@echo "Running NodeJS integration test with Dagger..."
	@go run ./test/integration/cmd/nodejs/run.go

.PHONY: test-integration
test-integration:
	@echo "Running all integration tests with Dagger..."
	@go run ./test/integration/cmd/run.go

generate-docs:
	@echo "Generating documentation..."
	@go run ./docs/generate-commands.go
	@echo "Documentation generated successfully!"

generate-schema:
	@echo "Generating schema..."
	@go run ./schema/generate-schema.go
	@echo "Schema generated successfully!"

.PHONY: generate-api
generate-api:
	@echo "Generating API clients from OpenAPI specs..."
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
		--config api/v1/push-codegen.yaml \
		api/v1/push.yaml > internal/api/client/push_client.gen.go
	@echo "API clients generated successfully!"

.PHONY: generate
generate: generate-api generate-docs generate-schema
	@echo "All code generation completed successfully!"

.PHONY: fmt
fmt:
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "Code formatted successfully!"