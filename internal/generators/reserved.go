package generators

import (
	"fmt"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/logger"
)

// FilterReservedFlags returns a copy of fs with any flag removed whose key,
// once transformed by the generator, collides with a symbol the generator
// reserves for itself (such as the exposed underlying OpenFeature client). A
// warning is emitted for each excluded flag so users understand why it is
// missing from the generated output. The reserved symbol always wins the
// collision.
//
// transform maps a flag key to the symbol name the generator would emit for it
// (e.g. ToPascal for Go). reserved holds the symbol names the generator
// reserves. generatorName is used only in the warning message.
func FilterReservedFlags(fs *flagset.Flagset, generatorName string, reserved map[string]bool, transform func(string) string) *flagset.Flagset {
	filtered := &flagset.Flagset{}
	for _, flag := range fs.Flags {
		symbol := transform(flag.Key)
		if reserved[symbol] {
			logger.Default.Warning(fmt.Sprintf(
				"Flag %q transforms to %q which is a reserved symbol in the %s generator. This flag will be excluded from the generated output.",
				flag.Key, symbol, generatorName,
			))
			continue
		}
		filtered.Flags = append(filtered.Flags, flag)
	}
	return filtered
}
