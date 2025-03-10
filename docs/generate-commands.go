package main

import (
	"fmt"
	"os"

	"github.com/open-feature/cli/cmd"
	"github.com/spf13/cobra/doc"
)

const docPath = "./docs/commands"

// Generates cobra docs of the cmd
func main() {
	linkHandler := func(name string) string {
		return name
	}

	filePrepender := func(filename string) string {
		return "<!-- markdownlint-disable-file -->\n<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->\n"
	}

	if err := doc.GenMarkdownTreeCustom(cmd.GetRootCmd(), docPath, filePrepender, linkHandler); err != nil {
		fmt.Fprintf(os.Stderr, "error generating docs: %v\n", err)
		os.Exit(1)
	}
}
