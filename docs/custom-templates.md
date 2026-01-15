# Custom Templates

The OpenFeature CLI supports custom templates that allow you to customize the generated code output. This is useful when you need to:

- Modify the structure of generated code to fit your project conventions
- Add custom imports, utilities, or wrappers
- Change naming conventions or code style
- Generate code for frameworks or libraries not yet supported

## Usage

Use the `--template` flag with any generate subcommand:

```bash
openfeature generate go --template ./my-custom-template.tmpl
openfeature generate react --template ./custom-react.tmpl
openfeature generate nodejs --template ./custom-nodejs.tmpl
```

## Getting Started

The easiest way to create a custom template is to start from an existing one:

1. Copy the default template for your target language from the CLI source:
   - Go: `internal/generators/golang/golang.tmpl`
   - React: `internal/generators/react/react.tmpl`
   - Node.js: `internal/generators/nodejs/nodejs.tmpl`
   - Python: `internal/generators/python/python.tmpl`
   - C#: `internal/generators/csharp/csharp.tmpl`
   - Java: `internal/generators/java/java.tmpl`
   - NestJS: `internal/generators/nestjs/nestjs.tmpl`

2. Modify the template to suit your needs

3. Use the `--template` flag to generate code with your custom template

## Template Syntax

Custom templates use Go's [text/template](https://pkg.go.dev/text/template) package. Refer to the official documentation for the full syntax, including conditionals, loops, and pipelines.

## Template Data

Templates have access to the following data structure:

```go
type TemplateData struct {
    Flagset struct {
        Flags []Flag
    }
    Params struct {
        OutputPath string
        Custom     any  // Language-specific parameters
    }
}

type Flag struct {
    Key          string   // The flag key (e.g., "enable-feature")
    Type         FlagType // The flag type (boolean, string, integer, float, object)
    Description  string   // Optional description of the flag
    DefaultValue any      // The default value for the flag
}
```

### Language-Specific Parameters

Some generators provide additional parameters in `.Params.Custom`:

**Go:**
- `.Params.Custom.GoPackage` - The Go package name
- `.Params.Custom.CLIVersion` - The CLI version used for generation

**C#:**
- `.Params.Custom.Namespace` - The C# namespace

**Java:**
- `.Params.Custom.JavaPackage` - The Java package name

## Template Functions

### Common Functions (Available in All Templates)

These functions are available in all templates:

| Function | Description | Example |
|----------|-------------|---------|
| `ToPascal` | Convert to PascalCase | `{{ .Key \| ToPascal }}` → `EnableFeature` |
| `ToCamel` | Convert to camelCase | `{{ .Key \| ToCamel }}` → `enableFeature` |
| `ToKebab` | Convert to kebab-case | `{{ .Key \| ToKebab }}` → `enable-feature` |
| `ToScreamingKebab` | Convert to SCREAMING-KEBAB-CASE | `{{ .Key \| ToScreamingKebab }}` → `ENABLE-FEATURE` |
| `ToSnake` | Convert to snake_case | `{{ .Key \| ToSnake }}` → `enable_feature` |
| `ToScreamingSnake` | Convert to SCREAMING_SNAKE_CASE | `{{ .Key \| ToScreamingSnake }}` → `ENABLE_FEATURE` |
| `ToUpper` | Convert to UPPERCASE | `{{ .Key \| ToUpper }}` → `ENABLE-FEATURE` |
| `ToLower` | Convert to lowercase | `{{ .Key \| ToLower }}` → `enable-feature` |
| `Quote` | Add double quotes | `{{ .Key \| Quote }}` → `"enable-feature"` |
| `QuoteString` | Quote if string type | `{{ .DefaultValue \| QuoteString }}` |

### Go-Specific Functions

| Function | Description |
|----------|-------------|
| `OpenFeatureType` | Convert flag type to OpenFeature method name (`Boolean`, `String`, `Int`, `Float`, `Object`) |
| `TypeString` | Convert flag type to Go type (`bool`, `string`, `int64`, `float64`, `map[string]any`) |
| `SupportImports` | Generate required imports based on flags |
| `ToMapLiteral` | Convert object value to Go map literal |

### React/Node.js/NestJS-Specific Functions

| Function | Description |
|----------|-------------|
| `OpenFeatureType` | Convert flag type to TypeScript type (`boolean`, `string`, `number`, `object`) |
| `ToJSONString` | Convert value to JSON string |

### Python-Specific Functions

| Function | Description |
|----------|-------------|
| `OpenFeatureType` | Convert flag type to Python type (`bool`, `str`, `int`, `float`, `object`) |
| `TypedGetMethodSync` | Get synchronous getter method name |
| `TypedGetMethodAsync` | Get async getter method name |
| `TypedDetailsMethodSync` | Get synchronous details method name |
| `TypedDetailsMethodAsync` | Get async details method name |
| `PythonBoolLiteral` | Convert boolean to Python literal (`True`/`False`) |
| `ToPythonDict` | Convert object value to Python dict literal |

### C#-Specific Functions

| Function | Description |
|----------|-------------|
| `OpenFeatureType` | Convert flag type to C# type (`bool`, `string`, `int`, `double`, `object`) |
| `FormatDefaultValue` | Format default value for C# |
| `ToCSharpDict` | Convert object value to C# dictionary literal |

### Java-Specific Functions

| Function | Description |
|----------|-------------|
| `OpenFeatureType` | Convert flag type to Java type (`Boolean`, `String`, `Integer`, `Double`, `Object`) |
| `FormatDefaultValue` | Format default value for Java |
| `ToMapLiteral` | Convert object value to Java Map literal |

## Example: Simple Go Template

Here's a minimal example of a custom Go template:

```go
// Code generated by OpenFeature CLI with custom template
package {{ .Params.Custom.GoPackage }}

import (
    "context"
    "github.com/open-feature/go-sdk/openfeature"
)

var client = openfeature.NewDefaultClient()

{{- range .Flagset.Flags }}
// Get{{ .Key | ToPascal }} returns the value of the "{{ .Key }}" flag.
// {{ if .Description }}{{ .Description }}{{ end }}
func Get{{ .Key | ToPascal }}(ctx context.Context, evalCtx openfeature.EvaluationContext) {{ .Type | TypeString }} {
    return client.{{ .Type | OpenFeatureType }}(ctx, {{ .Key | Quote }}, {{ .DefaultValue | QuoteString }}, evalCtx)
}
{{- end }}
```

## Example: Custom React Template

Here's an example that generates simple hooks without suspense:

```typescript
import { useFlag } from "@openfeature/react-sdk";

{{ range .Flagset.Flags }}
/**
 * {{ if .Description }}{{ .Description }}{{ else }}Feature flag{{ end }}
 * Default: {{ .DefaultValue }}
 */
export const use{{ .Key | ToPascal }} = () => {
  return useFlag({{ .Key | Quote }}, {{ .DefaultValue | QuoteString }});
};
{{ end }}
```

## Tips

1. **Test incrementally**: Make small changes and test the output frequently
2. **Use `--output` flag**: Direct output to a test directory while developing your template
3. **Preserve formatting**: The generators apply language-specific formatters after template execution (e.g., `gofmt` for Go)
4. **Handle edge cases**: Consider empty flag lists, missing descriptions, and different flag types
5. **Check the source**: Review the default templates in the CLI source for comprehensive examples
