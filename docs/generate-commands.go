package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/open-feature/cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const docPath = "./docs/commands"

// addStabilityToMarkdown adds stability information to the generated markdown content
func addStabilityToMarkdown(cmd *cobra.Command, content string) string {
	if stability, ok := cmd.Annotations["stability"]; ok {
		// Look for existing stability info to replace
		oldStabilityPattern := "\n> \\*\\*Stability\\*\\*: [a-z]+\n\n"
		hasExistingStability := regexp.MustCompile(oldStabilityPattern).MatchString(content)

		if hasExistingStability {
			// Replace existing stability info with the current one
			return regexp.MustCompile(oldStabilityPattern).ReplaceAllString(
				content,
				fmt.Sprintf("\n> **Stability**: %s\n\n", stability),
			)
		}

		// If no existing stability info, insert it
		// Look for the pattern of command title, description, and then either ### Synopsis or ```
		cmdNameLine := fmt.Sprintf("## openfeature%s", cmd.CommandPath()[11:])
		cmdNameIndex := strings.Index(content, cmdNameLine)

		if cmdNameIndex != -1 {
			// Find the end of the description section
			var insertPoint int
			synopsisIndex := strings.Index(content, "### Synopsis")
			codeBlockIndex := strings.Index(content, "```")

			if synopsisIndex != -1 {
				// If there's a Synopsis section, insert before it
				insertPoint = synopsisIndex
			} else if codeBlockIndex != -1 {
				// If there's a code block, insert before it
				insertPoint = codeBlockIndex
			} else {
				// Default to inserting after the description
				descStart := cmdNameIndex + len(cmdNameLine)
				nextNewline := strings.Index(content[descStart:], "\n\n")
				if nextNewline != -1 {
					insertPoint = descStart + nextNewline + 1
				} else {
					// Fallback to end of file
					insertPoint = len(content)
				}
			}

			stabilityInfo := fmt.Sprintf("\n> **Stability**: %s\n\n", stability)
			return content[:insertPoint] + stabilityInfo + content[insertPoint:]
		}
	}

	// If no stability annotation or couldn't find insertion point, return content unchanged
	return content
}

// Generates cobra docs of the cmd
func main() {
	linkHandler := func(name string) string {
		return name
	}

	filePrepender := func(filename string) string {
		return "<!-- markdownlint-disable-file -->\n<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->\n"
	}

	// Generate the markdown documentation
	if err := doc.GenMarkdownTreeCustom(cmd.GetRootCmd(), docPath, filePrepender, linkHandler); err != nil {
		fmt.Fprintf(os.Stderr, "error generating docs: %v\n", err)
		os.Exit(1)
	}

	// Apply the content modifier to all generated files
	// This is needed because Cobra doesn't expose a way to modify content during generation
	applyContentModifierToFiles(cmd.GetRootCmd(), docPath)
}

// applyContentModifierToFiles applies our content modifier to all generated markdown files
func applyContentModifierToFiles(root *cobra.Command, docPath string) {
	// Process the root command
	processCommandFile(root, fmt.Sprintf("%s/%s.md", docPath, root.Name()))

	// Process all descendant commands recursively
	processCommandTree(root, docPath)
}

// processCommandFile applies the content modifier to a single command's markdown file
func processCommandFile(cmd *cobra.Command, filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", filePath, err)
		return
	}

	// Apply our content modifier
	modifiedContent := addStabilityToMarkdown(cmd, string(content))

	// Only write the file if content was modified
	if modifiedContent != string(content) {
		err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing file %s: %v\n", filePath, err)
		}
	}
}

// processCommandTree recursively processes all commands in the command tree
func processCommandTree(cmd *cobra.Command, docPath string) {
	for _, subCmd := range cmd.Commands() {
		if !subCmd.IsAvailableCommand() || subCmd.IsAdditionalHelpTopicCommand() {
			continue
		}

		// Calculate the filename for this command
		fileName := getMarkdownFilename(cmd, subCmd)
		filePath := fmt.Sprintf("%s/%s", docPath, fileName)

		// Process this command's file
		processCommandFile(subCmd, filePath)

		// Process its children
		processCommandTree(subCmd, docPath)
	}
}

// getMarkdownFilename determines the markdown filename for a command based on its path
func getMarkdownFilename(parent *cobra.Command, cmd *cobra.Command) string {
	if parent.Name() == "openfeature" {
		return fmt.Sprintf("openfeature_%s.md", cmd.Name())
	}
	return fmt.Sprintf("openfeature_%s_%s.md", parent.Name(), cmd.Name())
}
