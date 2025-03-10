// Package filesystem contains the filesystem interface.
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var ViperKey = "filesystem"

// Get the filesystem interface from the viper configuration.
// If the filesystem interface is not set, the default filesystem interface is returned.
func FileSystem() afero.Fs {
	return viper.Get(ViperKey).(afero.Fs)
}

// Writes data to a file at the given path using the filesystem interface.
// If the file does not exist, it will be created, including all necessary directories.
// If the file exists, it will be overwritten.
func WriteFile(path string, data []byte) error {
	fs := FileSystem()
	if err := fs.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	f, err := fs.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file %q: %v", path, err)
	}
	defer f.Close()
	writtenBytes, err := f.Write(data)
	if err != nil {
		return fmt.Errorf("error writing contents to file %q: %v", path, err)
	}
	if writtenBytes != len(data) {
		return fmt.Errorf("error writing entire file %v: writtenBytes != expectedWrittenBytes", path)
	}

	return nil
}

// Checks if a file exists at the given path using the filesystem interface.
func Exists(path string) (bool, error) {
	fs := FileSystem()
	_, err := fs.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func init() {
	viper.SetDefault(ViperKey, afero.NewOsFs())
}
