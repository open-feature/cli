package java

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
	"github.com/spf13/afero"
)

func objectFlagset() *flagset.Flagset {
	return &flagset.Flagset{Flags: []flagset.Flag{
		{
			Key:         "retry_policy",
			Type:        flagset.ObjectType,
			Description: "Retry policy",
			DefaultValue: map[string]any{
				"backoff_seconds": 2.0,
				"max_attempts":    3.0,
				"ratio":           0.5,
				"big":             3000000000.0,
				"enabled":         true,
				"name":            "line1\nline2 \"q\" \\ é 😀",
				"nothing":         nil,
				"comment":         "end */ and \\u000a",
				"tags":            []any{"a", 1.0, map[string]any{"k": "v"}, []any{}},
				"nested":          map[string]any{"inner": map[string]any{"deep": false}, "empty": map[string]any{}},
				"k01":             1.0, "k02": 1.0, "k03": 1.0, "k04": 1.0, "k05": 1.0, "k06": 1.0, "k07": 1.0, "k08": 1.0,
			},
		},
		{Key: "plain", Type: flagset.BoolType, DefaultValue: true},
	}}
}

func generate(t *testing.T) string {
	t.Helper()
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)
	g := NewGenerator(objectFlagset())
	err := g.Generate(&generators.Params[Params]{
		OutputPath: "out",
		Custom:     Params{JavaPackage: "com.example.openfeature"},
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := afero.ReadFile(fs, filepath.Join("out", "OpenFeature.java"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestObjectFlagUsesValueType(t *testing.T) {
	out := generate(t)

	for _, want := range []string{
		"import dev.openfeature.sdk.Value;",
		"import dev.openfeature.sdk.MutableStructure;",
		"import java.util.List;",
		"Value retryPolicy(EvaluationContext ctx);",
		"FlagEvaluationDetails<Value> retryPolicyDetails(EvaluationContext ctx);",
		"public Value retryPolicy(EvaluationContext ctx) {",
		`client.getObjectValue("retry_policy", new Value(new MutableStructure()`,
		`client.getObjectDetails("retry_policy", new Value(new MutableStructure()`,
		`.add("nothing", new Value())`,
		`.add("big", new Value(3000000000.0))`,
		`.add("max_attempts", new Value(3))`,
		`new Value(List.<Value>of(`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("generated Java does not contain %q", want)
		}
	}
	for _, bad := range []string{"Map.of(", "Object retryPolicy", "FlagEvaluationDetails<Object>"} {
		if strings.Contains(out, bad) {
			t.Errorf("generated Java unexpectedly contains %q", bad)
		}
	}
}

// TestGeneratedJavaCompiles compiles the generated file with javac against the
// OpenFeature Java SDK. It is skipped unless javac and the SDK jar are
// available: set OPENFEATURE_JAVA_SDK_JAR to the path of dev.openfeature:sdk
// (and optionally SLF4J on the classpath via OPENFEATURE_JAVA_CLASSPATH_EXTRA).
func TestGeneratedJavaCompiles(t *testing.T) {
	javac, err := exec.LookPath("javac")
	if err != nil {
		t.Skip("javac not available")
	}
	jar := os.Getenv("OPENFEATURE_JAVA_SDK_JAR")
	if jar == "" {
		t.Skip("OPENFEATURE_JAVA_SDK_JAR not set")
	}
	out := generate(t)

	dir := t.TempDir()
	src := filepath.Join(dir, "OpenFeature.java")
	if err := os.WriteFile(src, []byte(out), 0o600); err != nil {
		t.Fatal(err)
	}
	cp := jar
	if extra := os.Getenv("OPENFEATURE_JAVA_CLASSPATH_EXTRA"); extra != "" {
		cp += string(os.PathListSeparator) + extra
	}
	cmd := exec.Command(javac, "-d", filepath.Join(dir, "classes"), "-cp", cp, src)
	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("javac failed: %v\n%s", err, combined)
	}
}

func TestToValueLiteral(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{
		{"nil", nil, "new Value()"},
		{"int", 3.0, "new Value(3)"},
		{"negative", -7.0, "new Value(-7)"},
		{"float", 0.15, "new Value(0.15)"},
		{"large integral", 1e10, "new Value(10000000000.0)"},
		{"string", "a\"b", `new Value("a\"b")`},
		{"empty list", []any{}, "new Value(List.<Value>of())"},
		{"map", map[string]any{"b": 1.0, "a": nil}, `new Value(new MutableStructure().add("a", new Value()).add("b", new Value(1)))`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toValueLiteral(tt.in); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
