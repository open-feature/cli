
.PHONY: test
test: 
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests passed successfully!"

generate-docs:
	@echo "Generating documentation..."
	@go run ./docs/generate-commands.go
	@echo "Documentation generated successfully!"

generate-schema:
	@echo "Generating schema..."
	@go run ./schema/generate-schema.go
	@echo "Schema generated successfully!"