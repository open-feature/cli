package generators

import (
	"github.com/open-feature/cli/internal/manifest"
)

// Represents the stability level of a generator
type Stability string

const (
	Unknown Stability = "unknown"
	Alpha   Stability = "alpha"
	Beta    Stability = "beta"
	Stable  Stability = "stable"
)

type Config struct {
	OutputPath string
}

type Generator interface {
	GetName() string
	GetStability() Stability
	SupportedFlagTypes() map[manifest.FlagType]bool
	Generate(config Config) error
}
