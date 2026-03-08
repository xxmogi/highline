package segment

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadBranch_NormalBranch(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := readBranch(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "main" {
		t.Errorf("expected 'main', got %q", got)
	}
}

func TestReadBranch_DetachedHead(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	sha := "abc1234def5678"
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(sha+"\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := readBranch(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != sha[:7] {
		t.Errorf("expected %q, got %q", sha[:7], got)
	}
}

func TestReadBranch_FeatureBranch(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "ref: refs/heads/feature/my-feature\n"
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := readBranch(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "feature/my-feature" {
		t.Errorf("expected 'feature/my-feature', got %q", got)
	}
}

func TestFindGitRoot_NotInRepo(t *testing.T) {
	// Change to a temp dir with no .git.
	tmp := t.TempDir()
	original, _ := os.Getwd()
	defer os.Chdir(original) //nolint

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	_, err := findGitRoot()
	if err == nil {
		t.Error("expected error when not in a git repository")
	}
}

func TestFindGitRoot_InRepo(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Change into a subdirectory.
	sub := filepath.Join(root, "sub", "dir")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	original, _ := os.Getwd()
	defer os.Chdir(original) //nolint

	if err := os.Chdir(sub); err != nil {
		t.Fatal(err)
	}

	got, err := findGitRoot()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != root {
		t.Errorf("expected root %q, got %q", root, got)
	}
}

func TestGitSegment_Name(t *testing.T) {
	g := NewGit(GitConfig{})
	if g.Name() != "git" {
		t.Errorf("expected 'git', got %q", g.Name())
	}
}

func TestGitSegment_Render_OutsideRepo(t *testing.T) {
	tmp := t.TempDir()
	original, _ := os.Getwd()
	defer os.Chdir(original) //nolint

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	g := NewGit(GitConfig{})
	_, _, _, err := g.Render()
	if err == nil {
		t.Error("expected error when outside git repo")
	}
}

func TestGitSegment_Render_CustomColors(t *testing.T) {
	root := t.TempDir()
	gitDir := filepath.Join(root, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	original, _ := os.Getwd()
	defer os.Chdir(original) //nolint
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	fg := 2
	bg := 5
	g := NewGit(GitConfig{FG: &fg, BG: &bg})
	_, gotFG, gotBG, err := g.Render()
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
