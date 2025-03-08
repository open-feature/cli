TODO: Add contributing guidelines

## Templates

### Data

The `TemplateData` struct is used to pass data to the templates.

### Built-in template functions

The following functions are automatically included in the templates:

#### ToPascal

Converts a string to `PascalCase`

```go
{{ "hello world" | ToPascal }} // HelloWorld
```

#### ToCamel

Converts a string to `camelCase`

```go
{{ "hello world" | ToCamel }} // helloWorld
```

#### ToKebab

Converts a string to `kebab-case`

```go
{{ "hello world" | ToKebab }} // hello-world
```

#### ToSnake

Converts a string to `snake_case`

```go
{{ "hello world" | ToSnake }} // hello_world
```

#### ToScreamingSnake

Converts a string to `SCREAMING_SNAKE_CASE`

```go
{{ "hello world" | ToScreamingSnake }} // HELLO_WORLD
```

#### ToUpper

Converts a string to `UPPER CASE`

```go
{{ "hello world" | ToUpper }} // HELLO WORLD
```

#### ToLower

Converts a string to `lower case`

```go
{{ "HELLO WORLD" | ToLower }} // hello world
```

#### ToTitle

Converts a string to `Title Case`

```go
{{ "hello world" | ToTitle }} // Hello World
```

#### Quote

Wraps a string in double quotes

```go
{{ "hello world" | Quote }} // "hello world"
```

#### QuoteString

Wraps only strings in double quotes

```go
{{ "hello world" | QuoteString }} // "hello world"
{{ 123 | QuoteString }} // 123
```

### Custom template functions

You can add custom template functions by passing a `FuncMap` to the `GenerateFile` function.
