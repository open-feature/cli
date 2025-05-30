package flagset

import (
	"strings"
	"testing"

	"github.com/open-feature/cli/internal/manifest"
)

// Sample test for FormatValidationError
func TestFormatValidationError_SortsByPath(t *testing.T) {
	issues := []manifest.ValidationError{
		{Path: "zeta.flag", Type: "boolean", Message: "must not be empty"},
		{Path: "alpha.flag", Type: "string", Message: "invalid value"},
		{Path: "beta.flag", Type: "number", Message: "must be greater than zero"},
	}

	output := FormatValidationError(issues)

	// The output should mention 'alpha.flag' before 'beta.flag', and 'beta.flag' before 'zeta.flag'
	alphaIdx := strings.Index(output, "flagPath: alpha.flag")
	betaIdx := strings.Index(output, "flagPath: beta.flag")
	zetaIdx := strings.Index(output, "flagPath: zeta.flag")

	if !(alphaIdx < betaIdx && betaIdx < zetaIdx) {
		t.Errorf("flag paths are not sorted: alphaIdx=%d, betaIdx=%d, zetaIdx=%d\nOutput:\n%s",
			alphaIdx, betaIdx, zetaIdx, output)
	}
}
