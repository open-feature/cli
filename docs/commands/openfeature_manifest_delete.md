<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## openfeature manifest delete

Delete a flag from the manifest

### Synopsis

Delete a flag from the manifest file by its key.

Examples:
  # Delete a flag named 'old-feature'
  openfeature manifest delete old-feature

  # Delete a flag from a specific manifest file
  openfeature manifest delete old-feature --manifest path/to/flags.json

```
openfeature manifest delete <flag-name> [flags]
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
      --debug             Enable debug logging
  -m, --manifest string   Path to the flag manifest (default "flags.json")
      --no-input          Disable interactive prompts
```

### SEE ALSO

* [openfeature manifest](openfeature_manifest.md)	 - Manage flag manifest files

