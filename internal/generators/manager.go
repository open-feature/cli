package generators

import (
	"sort"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// GeneratorCreator is a function that creates a generator command
type GeneratorCreator func() *cobra.Command

// GeneratorInfo contains metadata about a generator
type GeneratorInfo struct {
	Name             string
	Description      string
	Stability        Stability
	Creator          GeneratorCreator
}

// GeneratorManager maintains a registry of available generators
type GeneratorManager struct {
	generators map[string]GeneratorInfo
}

// NewGeneratorManager creates a new generator manager
func NewGeneratorManager() *GeneratorManager {
	return &GeneratorManager{
		generators: make(map[string]GeneratorInfo),
	}
}

// Register adds a generator to the registry
func (m *GeneratorManager) Register(cmdCreator func() *cobra.Command) {
	cmd := cmdCreator()
	m.generators[cmd.Use] = GeneratorInfo{
		Name:             cmd.Use,
		Description:      cmd.Short,
		Stability:        Stability(cmd.Annotations["stability"]),
		Creator:          cmdCreator,
	}
}

// GetAll returns all registered generators
func (m *GeneratorManager) GetAll() map[string]GeneratorInfo {
	return m.generators
}

// GetCommands returns cobra commands for all registered generators
func (m *GeneratorManager) GetCommands() []*cobra.Command {
	var commands []*cobra.Command
	
	for _, info := range m.generators {
		commands = append(commands, info.Creator())
	}
	
	return commands
}

// PrintGeneratorsTable prints a table of all available generators with their stability
func (m *GeneratorManager) PrintGeneratorsTable() error {
	tableData := [][]string{
		{"Generator", "Description", "Stability"},
	}
	
	// Get generator names for consistent ordering
	var names []string
	for name := range m.generators {
		names = append(names, name)
	}
	sort.Strings(names)
	
	for _, name := range names {
		info := m.generators[name]
		tableData = append(tableData, []string{
			name,
			info.Description,
			string(info.Stability),
		})
	}
	
	return pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

// DefaultManager is the default instance of the generator manager
var DefaultManager = NewGeneratorManager()
