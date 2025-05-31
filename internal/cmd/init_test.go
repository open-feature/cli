package cmd

import (
	"testing"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/spf13/afero"
)

func TestInitCmd(t *testing.T) {
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)
	outputFile := "flags-test.json"
	cmd := GetInitCmd()
	// global flag exists on root only.
	config.AddRootFlags(cmd)

	cmd.SetArgs([]string{
		"-m",
		outputFile,
	})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}
	compareOutput(t, "testdata/success_init.golden", outputFile, fs)
}
