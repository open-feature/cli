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
	name              string // test case name
	command           string // generator to run
	manifestGolden    string // path to the golden manifest file
	outputGolden      string // path to the golden output file
	outputPath        string // output directory (optional, defaults to "output")
	outputFile        string // output file name
	packageName       string // optional, used for Go (package-name), Java (package-name) and C# (namespace)
	templateFile      string // optional, path to a custom template file
	runtimeValidation *bool  // optional, defaults to true if nil
}

func TestGenerate(t *testing.T) {
	testCases := []generateTestCase{
		{
			name:           "Angular generation success",
			command:        "angular",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_angular.golden",
			outputFile:     "openfeature.generated.ts",
		},
		{
			name:           "Go generation success",
			command:        "go",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_go.golden",
			outputFile:     "testpackage_gen.go",
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
			name:           "NestJS generation success",
			command:        "nestjs",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_nestjs.golden",
			outputFile:     "openfeature-decorators.ts",
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
		{
			name:           "Java generation success",
			command:        "java",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/success_java.golden",
			outputFile:     "OpenFeature.java",
			packageName:    "com.example.openfeature",
		},
		{
			name:           "Angular generation with custom template",
			command:        "angular",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/custom_template/custom_angular.golden",
			outputFile:     "openfeature.generated.ts",
			templateFile:   "testdata/custom_template/custom_angular.tmpl",
		},
		{
			name:           "Go generation with custom template",
			command:        "go",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/custom_template/custom_go.golden",
			outputFile:     "testpackage_gen.go",
			packageName:    "testpackage",
			templateFile:   "testdata/custom_template/custom_go.tmpl",
		},
		{
			name:           "React generation with custom template",
			command:        "react",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/custom_template/custom_react.golden",
			outputFile:     "openfeature.ts",
			templateFile:   "testdata/custom_template/custom_react.tmpl",
		},
		{
			name:           "NodeJS generation with custom template",
			command:        "nodejs",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/custom_template/custom_nodejs.golden",
			outputFile:     "openfeature.ts",
			templateFile:   "testdata/custom_template/custom_nodejs.tmpl",
		},
		{
			name:           "Python generation with custom template",
			command:        "python",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/custom_template/custom_python.golden",
			outputFile:     "openfeature.py",
			templateFile:   "testdata/custom_template/custom_python.tmpl",
		},
		{
			name:           "CSharp generation with custom template",
			command:        "csharp",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/custom_template/custom_csharp.golden",
			outputFile:     "OpenFeature.g.cs",
			packageName:    "TestNamespace",
			templateFile:   "testdata/custom_template/custom_csharp.tmpl",
		},
		{
			name:           "Java generation with custom template",
			command:        "java",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/custom_template/custom_java.golden",
			outputFile:     "OpenFeature.java",
			packageName:    "com.example.openfeature",
			templateFile:   "testdata/custom_template/custom_java.tmpl",
		},
		{
			name:           "NestJS generation with custom template",
			command:        "nestjs",
			manifestGolden: "testdata/success_manifest.golden",
			outputGolden:   "testdata/custom_template/custom_nestjs.golden",
			outputFile:     "openfeature-decorators.ts",
			templateFile:   "testdata/custom_template/custom_nestjs.tmpl",
		},
		// Schema-typed object flag tests
		{
			name:           "Go generation with schema types",
			command:        "go",
			manifestGolden: "testdata/schema_manifest.golden",
			outputGolden:   "testdata/schema_go.golden",
			outputFile:     "testpackage_gen.go",
			packageName:    "testpackage",
		},
		{
			name:           "NodeJS generation with schema types",
			command:        "nodejs",
			manifestGolden: "testdata/schema_manifest.golden",
			outputGolden:   "testdata/schema_nodejs.golden",
			outputFile:     "openfeature.ts",
		},
		{
			name:           "React generation with schema types",
			command:        "react",
			manifestGolden: "testdata/schema_manifest.golden",
			outputGolden:   "testdata/schema_react.golden",
			outputFile:     "openfeature.ts",
		},
		{
			name:           "Angular generation with schema types",
			command:        "angular",
			manifestGolden: "testdata/schema_manifest.golden",
			outputGolden:   "testdata/schema_angular.golden",
			outputFile:     "openfeature.generated.ts",
		},
		{
			name:           "Python generation with schema types",
			command:        "python",
			manifestGolden: "testdata/schema_manifest.golden",
			outputGolden:   "testdata/schema_python.golden",
			outputFile:     "openfeature.py",
		},
		{
			name:           "CSharp generation with schema types",
			command:        "csharp",
			manifestGolden: "testdata/schema_manifest.golden",
			outputGolden:   "testdata/schema_csharp.golden",
			outputFile:     "OpenFeature.g.cs",
			packageName:    "TestNamespace",
		},
		{
			name:           "Java generation with schema types",
			command:        "java",
			manifestGolden: "testdata/schema_manifest.golden",
			outputGolden:   "testdata/schema_java.golden",
			outputFile:     "OpenFeature.java",
			packageName:    "com.example.openfeature",
		},
		{
			name:           "NestJS generation with schema types",
			command:        "nestjs",
			manifestGolden: "testdata/schema_manifest.golden",
			outputGolden:   "testdata/schema_nestjs.golden",
			outputFile:     "openfeature-decorators.ts",
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
			const memoryTemplatePath = "templates/custom.tmpl"

			// Use default output path if not specified
			outputPath := tc.outputPath
			if outputPath == "" {
				outputPath = "output"
			}

			// Prepare in-memory files
			fs := afero.NewMemMapFs()
			filesystem.SetFileSystem(fs)
			readOsFileAndWriteToMemMap(t, tc.manifestGolden, memoryManifestPath, fs)
			if tc.templateFile != "" {
				readOsFileAndWriteToMemMap(t, tc.templateFile, memoryTemplatePath, fs)
			}

			// Prepare command arguments
			args := []string{
				tc.command,
				"--manifest", memoryManifestPath,
				"--output", outputPath,
			}

			// Add parameters specific to each generator
			if tc.packageName != "" {
				switch tc.command {
				case "csharp":
					args = append(args, "--namespace", tc.packageName)
				case "go":
					args = append(args, "--package-name", tc.packageName)
				case "java":
					args = append(args, "--package-name", tc.packageName)
				}
			}

			// Add custom template flag if specified
			if tc.templateFile != "" {
				args = append(args, "--template", memoryTemplatePath)
			}

			// Add runtime-validation flag if specified
			if tc.runtimeValidation != nil {
				if *tc.runtimeValidation {
					args = append(args, "--runtime-validation")
				} else {
					args = append(args, "--runtime-validation=false")
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

// TestRuntimeValidationDisabled verifies that when --runtime-validation=false is set,
// generated code with schema-typed object flags includes type definitions but omits
// validation hooks. This ensures invalid provider objects won't be caught at runtime
// (compile-time safety only).
func TestRuntimeValidationDisabled(t *testing.T) {
	type runtimeValidationTestCase struct {
		name        string
		command     string
		outputFile  string
		packageName string
		// hookPatterns are strings that should be present when validation is ENABLED
		// and absent when validation is DISABLED
		hookPatterns []string
		// typePatterns are strings that should be present regardless of validation setting
		typePatterns []string
	}

	testCases := []runtimeValidationTestCase{
		{
			name:        "Go: hooks omitted when runtime-validation=false",
			command:     "go",
			outputFile:  "testpackage_gen.go",
			packageName: "testpackage",
			hookPatterns: []string{
				"themeCustomizationHook",
				"UnimplementedHook",
				"openfeature.WithHooks",
			},
			typePatterns: []string{
				"ThemeCustomizationValue",
				"PrimaryColor",
			},
		},
		{
			name:       "Node.js: hooks omitted when runtime-validation=false",
			command:    "nodejs",
			outputFile: "openfeature.ts",
			hookPatterns: []string{
				"createThemeCustomizationHook",
				"HookContext",
			},
			typePatterns: []string{
				"interface ThemeCustomization",
				"primaryColor: string",
			},
		},
		{
			name:       "React: hooks omitted when runtime-validation=false",
			command:    "react",
			outputFile: "openfeature.ts",
			hookPatterns: []string{
				"createThemeCustomizationHook",
				"HookContext",
			},
			typePatterns: []string{
				"interface ThemeCustomization",
				"primaryColor: string",
			},
		},
		{
			name:       "Python: hooks omitted when runtime-validation=false",
			command:    "python",
			outputFile: "openfeature.py",
			hookPatterns: []string{
				"ThemeCustomizationHook",
				"class ThemeCustomizationHook(Hook)",
			},
			typePatterns: []string{
				"class ThemeCustomization(TypedDict",
				"primaryColor: Required[str]",
			},
		},
		{
			name:        "C#: hooks omitted when runtime-validation=false",
			command:     "csharp",
			outputFile:  "OpenFeature.g.cs",
			packageName: "TestNamespace",
			hookPatterns: []string{
				"ThemeCustomizationHook",
				"class ThemeCustomizationHook : Hook",
			},
			typePatterns: []string{
				"record ThemeCustomization(",
				"string PrimaryColor",
			},
		},
		{
			name:        "Java: hooks omitted when runtime-validation=false",
			command:     "java",
			outputFile:  "OpenFeature.java",
			packageName: "com.example.openfeature",
			hookPatterns: []string{
				"ThemeCustomizationHook",
				"class ThemeCustomizationHook implements Hook",
			},
			typePatterns: []string{
				"record ThemeCustomization(",
				"String primaryColor",
			},
		},
	}

	runtimeValidationFalse := false

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := GetGenerateCmd()
			config.AddRootFlags(cmd)

			const memoryManifestPath = "manifest/path.json"
			outputPath := "output"

			fs := afero.NewMemMapFs()
			filesystem.SetFileSystem(fs)
			readOsFileAndWriteToMemMap(t, "testdata/schema_manifest.golden", memoryManifestPath, fs)

			args := []string{
				tc.command,
				"--manifest", memoryManifestPath,
				"--output", outputPath,
				"--runtime-validation=false",
			}

			if tc.packageName != "" {
				switch tc.command {
				case "csharp":
					args = append(args, "--namespace", tc.packageName)
				case "go":
					args = append(args, "--package-name", tc.packageName)
				case "java":
					args = append(args, "--package-name", tc.packageName)
				}
			}

			cmd.SetArgs(args)
			err := cmd.Execute()
			if err != nil {
				t.Fatalf("generation failed: %v", err)
			}

			got, err := afero.ReadFile(fs, filepath.Join(outputPath, tc.outputFile))
			if err != nil {
				t.Fatalf("error reading output file: %v", err)
			}
			output := string(got)

			// Type definitions should still be present
			for _, pattern := range tc.typePatterns {
				if !strings.Contains(output, pattern) {
					t.Errorf("expected type pattern %q to be present in output (types should exist even without runtime validation)", pattern)
				}
			}

			// Hook patterns should be absent
			for _, pattern := range tc.hookPatterns {
				if strings.Contains(output, pattern) {
					t.Errorf("expected hook pattern %q to be absent in output when runtime-validation=false", pattern)
				}
			}

			_ = runtimeValidationFalse // referenced for documentation
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

// normalizeLines trims trailing whitespace and carriage returns from each line.
// This helps ensure consistent comparison by ignoring formatting differences like indentation or line endings.
func normalizeLines(input []string) []string {
	normalized := make([]string, len(input))
	for i, line := range input {
		// Trim right whitespace and convert \r\n or \r to \n
		normalized[i] = strings.TrimRight(line, " \t\r")
	}
	return normalized
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
	wantLines := normalizeLines(strings.Split(string(want), "\n"))
	gotLines := normalizeLines(strings.Split(string(got), "\n"))

	if diff := cmp.Diff(wantLines, gotLines); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}
