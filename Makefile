.PHONY: test
test: 
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests passed successfully!"

# Dagger-based integration tests
.PHONY: test-csharp-dagger
test-csharp-dagger:
	@echo "Running C# integration test with Dagger..."
	@go run ./test/integration/cmd/csharp/run.go

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

.PHONY: fmt
fmt:
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "Code formatted successfully!"