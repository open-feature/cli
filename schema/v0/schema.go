// Package schema embeds the flag manifest into a code module.
package schema

import _ "embed"

// Schema contains the embedded flag manifest schema.
//
//go:embed flag_manifest.json
var SchemaFile string
