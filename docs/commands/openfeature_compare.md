<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## openfeature compare

Compare two feature flag manifests

### Synopsis

Compare two OpenFeature flag manifests and display the differences in a structured format.

By default, shows what HAS changed in the manifest compared to the target (receiving perspective).
Use --reverse to show what WILL change when the manifest is pushed to the target (sending perspective).

Examples:
  # Show what changed in local compared to main (default)
  openfeature compare --manifest local.json --against main.json

  # Preview what will change when pushing to remote
  openfeature compare --manifest local.json --against remote.json --reverse

```
openfeature compare [flags]
```

### Options

```
  -a, --against string       Path to the target manifest file to compare against
  -h, --help                 help for compare
  -i, --ignore stringArray   Field pattern to ignore during comparison (can be specified multiple times). Supports shorthand (e.g., 'description') and full paths with wildcards (e.g., 'flags.*.description', 'metadata.*')
  -o, --output string        Output format. Valid formats: tree, flat, json, yaml (default "tree")
      --reverse              Reverse comparison direction. Shows what WILL change when manifest is pushed to target (sending perspective) instead of what HAS changed in manifest compared to target (receiving perspective)
```

### Options inherited from parent commands

```
      --debug             Enable debug logging
  -m, --manifest string   Path to the flag manifest (default "flags.json")
      --no-input          Disable interactive prompts
```

### SEE ALSO

* [openfeature](openfeature.md)	 - CLI for OpenFeature.

