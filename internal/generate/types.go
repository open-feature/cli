package generate

import (
	"github.com/open-feature/cli/internal/manifest"
)

type Generator interface {
	// Generate generates the code for the given input.
	Generate(input Config) error
	// SupportedFlagTypes returns the flag types supported by the generator.
	SupportedFlagTypes() map[manifest.FlagType]bool
}