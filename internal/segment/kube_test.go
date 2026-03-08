package segment

import (
	"os"
	"path/filepath"
	"testing"
)

func writeKubeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCurrentKubeContext_Basic(t *testing.T) {
	config := `apiVersion: v1
kind: Config
current-context: my-cluster
contexts:
- name: my-cluster
`
	path := writeKubeConfig(t, config)
	t.Setenv("KUBECONFIG", path)

	got, err := currentKubeContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "my-cluster" {
		t.Errorf("expected 'my-cluster', got %q", got)
	}
}

func TestCurrentKubeContext_Quoted(t *testing.T) {
	config := `current-context: "prod-us-east"`
	path := writeKubeConfig(t, config)
	t.Setenv("KUBECONFIG", path)

	got, err := currentKubeContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "prod-us-east" {
		t.Errorf("expected 'prod-us-east', got %q", got)
	}
}

func TestCurrentKubeContext_Empty(t *testing.T) {
	config := `current-context: `
	path := writeKubeConfig(t, config)
	t.Setenv("KUBECONFIG", path)

	_, err := currentKubeContext()
	if err == nil {
		t.Error("expected error for empty context")
	}
}

func TestCurrentKubeContext_Null(t *testing.T) {
	config := `current-context: null`
	path := writeKubeConfig(t, config)
	t.Setenv("KUBECONFIG", path)

	_, err := currentKubeContext()
	if err == nil {
		t.Error("expected error for null context")
	}
}

func TestCurrentKubeContext_MissingFile(t *testing.T) {
	t.Setenv("KUBECONFIG", "/nonexistent/path/config")

	_, err := currentKubeContext()
	if err == nil {
		t.Error("expected error for missing kubeconfig")
	}
}

func TestCurrentKubeContext_NoContextField(t *testing.T) {
	config := `apiVersion: v1
kind: Config
`
	path := writeKubeConfig(t, config)
	t.Setenv("KUBECONFIG", path)

	_, err := currentKubeContext()
	if err == nil {
		t.Error("expected error when current-context field is absent")
	}
}

func TestKubeConfigPath_EnvVar(t *testing.T) {
	t.Setenv("KUBECONFIG", "/custom/kube/config")
	got := kubeConfigPath()
	if got != "/custom/kube/config" {
		t.Errorf("expected '/custom/kube/config', got %q", got)
	}
}

func TestKubeConfigPath_ColonSeparated(t *testing.T) {
	t.Setenv("KUBECONFIG", "/first/config:/second/config")
	got := kubeConfigPath()
	if got != "/first/config" {
		t.Errorf("expected '/first/config', got %q", got)
	}
}

func TestKubeSegment_Name(t *testing.T) {
	k := NewKube(KubeConfig{})
	if k.Name() != "kube" {
		t.Errorf("expected 'kube', got %q", k.Name())
	}
}

func TestKubeSegment_Render(t *testing.T) {
	config := `current-context: test-ctx`
	path := writeKubeConfig(t, config)
	t.Setenv("KUBECONFIG", path)

	k := NewKube(KubeConfig{})
	text, _, _, err := k.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text == "" {
		t.Error("expected non-empty text")
	}
	// Should contain the context name.
	if !contains(text, "test-ctx") {
		t.Errorf("expected 'test-ctx' in output, got %q", text)
	}
}

func TestKubeSegment_Render_CustomColors(t *testing.T) {
	config := `current-context: test-ctx`
	path := writeKubeConfig(t, config)
	t.Setenv("KUBECONFIG", path)

	fg := 1
	bg := 3
	k := NewKube(KubeConfig{FG: &fg, BG: &bg})
	_, gotFG, gotBG, err := k.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFG != fg {
		t.Errorf("expected fg=%d, got %d", fg, gotFG)
	}
	if gotBG != bg {
		t.Errorf("expected bg=%d, got %d", bg, gotBG)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
