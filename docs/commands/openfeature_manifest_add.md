<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## openfeature manifest add

Add a new flag to the manifest

### Synopsis

Add a new flag to the manifest file with the specified configuration.

Interactive Mode:
  When flags are omitted, the command prompts interactively for missing values:
  - Flag type (defaults to boolean if not specified)
  - Default value (required)
  - Description (optional, press Enter to skip)
  
  Use --no-input to disable interactive prompts (required for CI/automation).

Examples:
  # Interactive mode - prompts for type, value, and description
  openfeature manifest add new-feature

  # Add a boolean flag (default type)
  openfeature manifest add new-feature --default-value false

  # Add a string flag with description
  openfeature manifest add welcome-message --type string --default-value "Hello!" --description "Welcome message for users"

  # Add an integer flag
  openfeature manifest add max-retries --type integer --default-value 3

  # Add a float flag
  openfeature manifest add discount-rate --type float --default-value 0.15

  # Add an object flag
  openfeature manifest add config --type object --default-value '{"key":"value"}'
  
  # Disable interactive prompts (for automation)
  openfeature manifest add my-flag --default-value true --no-input

```
openfeature manifest add [flag-name] [flags]
```

### Options

```
  -d, --default-value string   Default value for the flag (required)
      --description string     Description of the flag
  -h, --help                   help for add
  -t, --type string            Type of the flag (boolean, string, integer, float, object) (default "boolean")
```

### Options inherited from parent commands

```
      --debug             Enable debug logging
  -m, --manifest string   Path to the flag manifest (default "flags.json")
      --no-input          Disable interactive prompts
```

### SEE ALSO

* [openfeature manifest](openfeature_manifest.md)	 - Manage flag manifest files

