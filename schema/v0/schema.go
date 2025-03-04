// Package schema embeds the flag manifest into a code module.
package schema

import _ "embed"

// Schema contains the embedded flag manifest schema.
//
//go:embed flag-manifest.json
var SchemaFile string
