.PHONY: test
test: 
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests passed successfully!"

.PHONY: test-csharp
test-csharp:
	@echo "Running C# integration test..."
	@./test/csharp-integration/test-compilation.sh

.PHONY: test-csharp-dagger
test-csharp-dagger:
	@echo "Running C# integration test with Dagger..."
	@cd test/csharp-integration && go run dagger.go

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