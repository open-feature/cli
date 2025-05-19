<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## openfeature compare

Compare two feature flag manifests

### Synopsis

Compare two OpenFeature flag manifests and display the differences in a structured format.

```
openfeature compare [flags]
```

### Options

```
  -a, --against string   Path to the target manifest file to compare against
  -h, --help             help for compare
  -o, --output string    Output format. Valid formats: tree, flat, json, yaml (default "tree")
```

### Options inherited from parent commands

```
      --debug             Enable debug logging
  -m, --manifest string   Path to the flag manifest (default "flags.json")
      --no-input          Disable interactive prompts
```

### SEE ALSO

* [openfeature](openfeature.md)	 - CLI for OpenFeature.

