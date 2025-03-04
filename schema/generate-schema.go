package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/open-feature/cli/internal/manifest"
)

const schemaPath = "schema/v0/flag-manifest.json"

func main() {
	schema := manifest.ToJSONSchema()
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to marshal JSON schema: %w", err))
	}

	if err := os.MkdirAll("schema/v0", os.ModePerm); err != nil {
		log.Fatal(fmt.Errorf("failed to create directory: %w", err))
	}

	file, err := os.Create(schemaPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create file: %w", err))
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		log.Fatal(fmt.Errorf("failed to write JSON schema to file: %w", err));
	}

	fmt.Println("JSON schema generated successfully at " + schemaPath)
}