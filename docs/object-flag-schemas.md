# Object Flag Schemas

The OpenFeature CLI supports an optional `schema` field on object-type flags that enables **type-safe code generation** and **runtime validation**. When a schema is provided, the generated code includes language-specific typed structures instead of generic object types, and optionally includes runtime validation hooks that verify provider responses at evaluation time.

## Overview

By default, object flags return generic types in generated code (`any` in Go, `JsonValue` in TypeScript, `object` in Python, etc.). This forces developers to manually cast values and provides no compile-time safety.

With the `schema` field, generated code includes:

- **Typed structures**: Language-native types (Go structs, TypeScript interfaces, Java records, C# records, Python TypedDicts)
- **Runtime validation hooks**: OpenFeature `after` hooks that validate the provider's response matches the expected shape
- **Compile-time safety**: IDE autocomplete, type checking, and refactoring support

## Schema Definition

The `schema` field uses a **subset of JSON Schema** to describe the shape of an object flag's value. Add it to any object-type flag in your manifest:

```json
{
  "flags": {
    "themeCustomization": {
      "flagType": "object",
      "defaultValue": {
        "primaryColor": "#007bff",
        "secondaryColor": "#6c757d",
        "header": {
          "fontSize": 16,
          "visible": true
        }
      },
      "description": "Allows customization of theme colors.",
      "schema": {
        "type": "object",
        "properties": {
          "primaryColor": { "type": "string" },
          "secondaryColor": { "type": "string" },
          "header": {
            "type": "object",
            "properties": {
              "fontSize": { "type": "integer" },
              "visible": { "type": "boolean" }
            },
            "required": ["fontSize"]
          }
        },
        "required": ["primaryColor"]
      }
    }
  }
}
```

### Supported JSON Schema Keywords

The schema field supports the following JSON Schema keywords:

| Keyword | Description | Applies To |
|---------|-------------|------------|
| `type` | **Required.** The data type. | All schemas |
| `properties` | Map of property names to their schemas. | `object` |
| `required` | Array of property names that must be present. | `object` |
| `items` | Schema for array elements. | `array` |
| `additionalProperties` | Whether extra properties are allowed (boolean). | `object` |

### Supported Types

| Type | Description | Language Mapping |
|------|-------------|-----------------|
| `string` | Text value | Go: `string`, TS: `string`, Java: `String`, C#: `string`, Python: `str` |
| `number` | Any numeric value | Go: `float64`, TS: `number`, Java: `Double`, C#: `double`, Python: `float` |
| `integer` | Whole number | Go: `int64`, TS: `number`, Java: `Integer`, C#: `int`, Python: `int` |
| `boolean` | True/false | Go: `bool`, TS: `boolean`, Java: `Boolean`, C#: `bool`, Python: `bool` |
| `object` | Nested object (recursive) | Go: inline struct, TS: inline interface, Java/C#: nested record, Python: nested TypedDict |
| `array` | List of items | Go: slice, TS: `T[]`, Java: `List<T>`, C#: `List<T>`, Python: `list[T]` |

### Nesting

Schemas support arbitrary nesting. Nested `object` types generate fully typed nested structures:

```json
{
  "schema": {
    "type": "object",
    "properties": {
      "header": {
        "type": "object",
        "properties": {
          "fontSize": { "type": "integer" },
          "visible": { "type": "boolean" }
        },
        "required": ["fontSize"]
      }
    }
  }
}
```

For languages that support inline anonymous types (Go, TypeScript), nesting is expressed inline. For languages that require named types (Java, C#, Python), the CLI generates depth-first named types with compound names (e.g., `ThemeCustomizationHeader` for a `header` property on a `ThemeCustomization` flag).

### Arrays

Array types use the `items` keyword to define the element schema:

```json
{
  "schema": {
    "type": "object",
    "properties": {
      "tags": {
        "type": "array",
        "items": { "type": "string" }
      },
      "limits": {
        "type": "array",
        "items": { "type": "number" }
      }
    }
  }
}
```

## Default Value Validation

When a `schema` is provided, the CLI validates the flag's `defaultValue` against the schema at manifest load time. This catches authoring mistakes early:

- Missing required properties
- Type mismatches (e.g., string where integer is expected)
- Additional properties when `additionalProperties: false`
- Nested validation (recursive)
- Array element type validation

If validation fails, the CLI reports errors with paths like `flags.myFlag.defaultValue.header.fontSize`.

## Runtime Validation

When code is generated with `--runtime-validation` (enabled by default), the CLI generates OpenFeature `after` hooks that validate the provider's response at evaluation time. This catches cases where a feature flag provider returns an object that doesn't match the expected schema.

### What the Hooks Validate

The generated hooks recursively validate the full schema:

1. **Object type check**: The resolved value is an object (map/dict/structure), not null, not an array, and not a primitive
2. **Required property check**: Each property listed in `schema.required` exists at every level of nesting
3. **Property type validation**: Each present property's value matches its declared type (`string`, `number`, `integer`, `boolean`)
4. **Nested object validation**: Nested objects are validated recursively with the same checks
5. **Array validation**: Array properties are type-checked, and each element is validated against `schema.items`
6. **Additional properties check**: When `additionalProperties: false`, any unexpected keys cause a validation error

If validation fails, the hook raises an error. The OpenFeature SDK specification defines that when an `after` hook errors, the SDK should return the default value instead. This means malformed provider responses gracefully fall back to the declared `defaultValue`.

### Disabling Runtime Validation

To generate code with compile-time types only (no validation hooks):

```bash
openfeature generate go --runtime-validation=false
openfeature generate react --runtime-validation=false
```

This is useful when:
- You trust the provider always returns correctly shaped objects
- You want to minimize generated code size
- You want to handle validation in your own application logic

### Per-Language Behavior

| Language | Compile-Time Types | Runtime Hooks | Notes |
|----------|-------------------|---------------|-------|
| Go | Structs with JSON tags | `After` hook on `UnimplementedHook` | Uses `json.Marshal`/`Unmarshal` round-trip for type conversion. Type name uses `Value` suffix (e.g., `ThemeCustomizationValue`) to avoid conflict with the generated variable. |
| Node.js | TypeScript interfaces | `Hook` with `after` callback | Hooks injected via spread into `FlagEvaluationOptions`. |
| React | TypeScript interfaces | `Hook` with `after` callback | Hooks injected into `useFlag`/`useSuspenseFlag` options. |
| Angular | TypeScript interfaces | None | Compile-time types only. Angular directive/service architecture doesn't support per-evaluation hooks. |
| NestJS | TypeScript interfaces (imported from Node.js file) | None | Compile-time types only. Decorators don't support per-evaluation hooks. Runtime validation is handled by the Node.js generated client. |
| Java | Records with `@Nullable` annotations | `Hook<Object>` with `after` method | Uses `ObjectMapper.convertValue()` for type conversion. Nested types are static inner records. |
| C# | Records | `Hook` with `AfterAsync` method | Uses `JsonSerializer.Deserialize` for type conversion. Nested types are nested records. |
| Python | TypedDicts with `Required[]` | `Hook` subclass with `after` method | Hooks injected via `FlagEvaluationOptions`. Nested types are separate TypedDict classes. |

## Backward Compatibility

The `schema` field is entirely optional:

- Manifests without `schema` on object flags continue to work unchanged
- Object flags without `schema` still generate generic types (`any`, `JsonValue`, `Value`, etc.)
- The `--runtime-validation` flag has no effect on flags without schemas

## Limitations

### JSON Schema Subset

The `schema` field supports a **limited subset** of JSON Schema. The following JSON Schema features are **not supported**:

- `enum` / `const` (string/number constraints)
- `pattern` (regex validation for strings)
- `minimum` / `maximum` / `exclusiveMinimum` / `exclusiveMaximum` (numeric ranges)
- `minLength` / `maxLength` (string length)
- `minItems` / `maxItems` / `uniqueItems` (array constraints)
- `minProperties` / `maxProperties` (object property count)
- `oneOf` / `anyOf` / `allOf` / `not` (composition)
- `$ref` / `$defs` (schema references, except internally for the manifest JSON Schema itself)
- `format` (e.g., `date-time`, `email`, `uri`)
- `default` (property-level defaults)
- `if` / `then` / `else` (conditional schemas)
- `patternProperties` (regex-based property schemas)
- `nullable` / type arrays (e.g., `"type": ["string", "null"]`)

The schema format is designed to be forward-compatible with future additions like `enum` support.

### Go SDK `ObjectValueDetails` Behavior

The Go SDK's `ObjectValueDetails` method has an inconsistency where the `Value` field in the returned details contains the resolved value (not the default) even when an `after` hook returns an error. This differs from all other typed `*ValueDetails` methods in the Go SDK. This is a known upstream issue and will be fixed in the Go SDK. No workaround is applied in the generated code.

### Angular and NestJS

Angular and NestJS generators produce compile-time types only. The Angular `FeatureFlagDirective` and NestJS decorator patterns don't support injecting per-evaluation hooks through the generated code. If you need runtime validation for Angular or NestJS, use the Node.js generated client with runtime validation enabled and consume its output.

## Example: Generated Code

Given the `themeCustomization` schema from the example above:

### Go

```go
type ThemeCustomizationValue struct {
    Header struct {
        FontSize int64 `json:"fontSize"`
        Visible  bool  `json:"visible,omitempty"`
    } `json:"header,omitempty"`
    PrimaryColor   string `json:"primaryColor"`
    SecondaryColor string `json:"secondaryColor,omitempty"`
}
```

### TypeScript (Node.js / React)

```typescript
export interface ThemeCustomization {
  header?: {
    fontSize: number;
    visible?: boolean;
  };
  primaryColor: string;
  secondaryColor?: string;
}
```

Run the generator for your target language to see the full typed output for Java, C#, and Python.
