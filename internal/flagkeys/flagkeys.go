// Package commonflags contains keys for all command-line flags related to openfeature CLI.
package flagkeys

import "github.com/spf13/viper"

const (
	// `generate` flags:
	// FlagManifestPath is the key for the flag that stores the flag manifest path.
	FlagManifestPath = "flag_manifest_path"
	// OutputPath is the key for the flag that stores the output path.
	OutputPath = "output_path"

	// `generate go` flags:
	// GoPackageName is the key for the flag that stores the Golang package name.
	GoPackageName = "package_name"

	// `generate java` flags:
	// JavaPackageName is the key for the flag that stores the Java package name.
	JavaPackageName = "package_name"

	//internal keys:
	// FileSystem is the key for the flag that stores the filesystem interface.
	FileSystem = "filesystem"
)

func init() {
	viper.SetDefault(FileSystem, "local")
}
