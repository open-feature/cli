package config

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAddCppGenerateFlags_Default(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddCppGenerateFlags(cmd)

	flag := cmd.Flags().Lookup(CppNamespaceName)
	if flag == nil {
		t.Fatalf("flag %q not found", CppNamespaceName)
	}

	if flag.DefValue != DefaultCppNamespace {
		t.Fatalf("flag default value = %q; want %q", flag.DefValue, DefaultCppNamespace)
	}

	got := GetCppNamespace(cmd)
	if got != DefaultCppNamespace {
		t.Fatalf("GetCppNamespace(default) = %q; want %q", got, DefaultCppNamespace)
	}
}

func TestGetCppNamespace_CustomValue(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddCppGenerateFlags(cmd)

	err := cmd.Flags().Set(CppNamespaceName, "custom_namespace")
	if err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	got := GetCppNamespace(cmd)
	if got != "custom_namespace" {
		t.Fatalf("GetCppNamespace(custom) = %q; want %q", got, "custom_namespace")
	}
}

func TestAddCSharpGenerateFlags_Default(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddCSharpGenerateFlags(cmd)

	flag := cmd.Flags().Lookup(CSharpNamespaceName)
	if flag == nil {
		t.Fatalf("flag %q not found", CSharpNamespaceName)
	}

	if flag.DefValue != DefaultCSharpNamespace {
		t.Fatalf("flag default value = %q; want %q", flag.DefValue, DefaultCSharpNamespace)
	}

	got := GetCSharpNamespace(cmd)
	if got != DefaultCSharpNamespace {
		t.Fatalf("GetCSharpNamespace(default) = %q; want %q", got, DefaultCSharpNamespace)
	}
}

func TestAddJavaGenerateFlags_Default(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddJavaGenerateFlags(cmd)

	flag := cmd.Flags().Lookup(JavaPackageFlagName)
	if flag == nil {
		t.Fatalf("flag %q not found", JavaPackageFlagName)
	}

	if flag.DefValue != DefaultJavaPackageName {
		t.Fatalf("flag default value = %q; want %q", flag.DefValue, DefaultJavaPackageName)
	}

	got := GetJavaPackageName(cmd)
	if got != DefaultJavaPackageName {
		t.Fatalf("GetJavaPackageName(default) = %q; want %q", got, DefaultJavaPackageName)
	}
}

func TestAddGoGenerateFlags_Default(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddGoGenerateFlags(cmd)

	flag := cmd.Flags().Lookup(GoPackageFlagName)
	if flag == nil {
		t.Fatalf("flag %q not found", GoPackageFlagName)
	}

	if flag.DefValue != DefaultGoPackageName {
		t.Fatalf("flag default value = %q; want %q", flag.DefValue, DefaultGoPackageName)
	}

	got := GetGoPackageName(cmd)
	if got != DefaultGoPackageName {
		t.Fatalf("GetGoPackageName(default) = %q; want %q", got, DefaultGoPackageName)
	}
}
