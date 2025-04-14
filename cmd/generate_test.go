package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"

	"github.com/spf13/afero"
)

// generateTestCase holds the configuration for each generate test
type generateTestCase struct {
	name           string // test case name
	command        string // generator to run
	manifestGolden string // path to the golden manifest file
	outputGolden   string // path to the golden output file
	outputPath     string // output directory (optional, defaults to "output")
	outputFile     string // output file name
	packageName    string // optional, used for Go (package-name) and C# (namespace)
}

func TestGenerate(t *testing.T) {
	testCases := []generateTestCase{
		{
			name:           "Go generation success",
			command:        "go",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_go.golden",
			outputFile:     "testpackage.go",
			packageName:    "testpackage",
		},
		{
			name:           "React generation success",
			command:        "react",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_react.golden",
			outputFile:     "openfeature.ts",
		},
		{
			name:           "NodeJS generation success",
			command:        "nodejs",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_nodejs.golden",
			outputFile:     "openfeature.ts",
		},
		{
			name:           "Python generation success",
			command:        "python",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_python.golden",
			outputFile:     "openfeature.py",
		},
		{
			name:           "CSharp generation success",
			command:        "csharp",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_csharp.golden",
			outputFile:     "OpenFeature.g.cs",
			packageName:    "TestNamespace", // Using packageName field for namespace
		},
		// Add more test cases here as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := GetGenerateCmd()

			// global flag exists on root only.
			config.AddRootFlags(cmd)

			// Constant paths
			const memoryManifestPath = "manifest/path.json"

			// Use default output path if not specified
			outputPath := tc.outputPath
			if outputPath == "" {
				outputPath = "output"
			}

			// Prepare in-memory files
			fs := afero.NewMemMapFs()
			filesystem.SetFileSystem(fs)
			readOsFileAndWriteToMemMap(t, tc.manifestGolden, memoryManifestPath, fs)

			// Prepare command arguments
			args := []string{
				tc.command,
				"--manifest", memoryManifestPath,
				"--output", outputPath,
			}

			// Add parameters specific to each generator
			if tc.packageName != "" {
				if tc.command == "csharp" {
					args = append(args, "--namespace", tc.packageName)
				} else if tc.command == "go" {
					args = append(args, "--package-name", tc.packageName)
				}
			}

			cmd.SetArgs(args)

			// Run command
			err := cmd.Execute()
			if err != nil {
				t.Error(err)
			}

			// Compare result
			compareOutput(t, tc.outputGolden, filepath.Join(outputPath, tc.outputFile), fs)
		})
	}
}

func readOsFileAndWriteToMemMap(t *testing.T, inputPath string, memPath string, memFs afero.Fs) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("error reading file %q: %v", inputPath, err)
	}
	if err := memFs.MkdirAll(filepath.Dir(memPath), os.ModePerm); err != nil {
		t.Fatalf("error creating directory %q: %v", filepath.Dir(memPath), err)
	}
	f, err := memFs.Create(memPath)
	if err != nil {
		t.Fatalf("error creating file %q: %v", memPath, err)
	}
	defer f.Close()
	writtenBytes, err := f.Write(data)
	if err != nil {
		t.Fatalf("error writing contents to file %q: %v", memPath, err)
	}
	if writtenBytes != len(data) {
		t.Fatalf("error writing entire file %v: writtenBytes != expectedWrittenBytes", memPath)
	}
}

func compareOutput(t *testing.T, testFile, memoryOutputPath string, fs afero.Fs) {
	want, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("error reading file %q: %v", testFile, err)
	}

	got, err := afero.ReadFile(fs, memoryOutputPath)
	if err != nil {
		t.Fatalf("error reading file %q: %v", memoryOutputPath, err)
	}

	// Convert to string arrays by splitting on newlines
	wantLines := strings.Split(string(want), "\n")
	gotLines := strings.Split(string(got), "\n")

	if diff := cmp.Diff(wantLines, gotLines); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}
