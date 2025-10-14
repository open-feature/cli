<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## openfeature manifest add

Add a new flag to the manifest

### Synopsis

Add a new flag to the manifest file with the specified configuration.

Examples:
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

