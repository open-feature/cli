// Package api provides generated API clients for OpenFeature CLI
package api

// Generate API clients from OpenAPI specifications
// Run: go generate ./internal/api/...

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config=../../api/v0/sync-codegen.yaml ../../api/v0/sync.yaml -o client/sync_client.gen.go
