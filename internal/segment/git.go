package segment

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/user/highline/internal/color"
)

// GitConfig holds optional color configuration for GitSegment.
type GitConfig struct {
	FG *int // nil = color.Black
	BG *int // nil = color.Yellow
}

// GitSegment renders the current Git branch and dirty status.
type GitSegment struct {
	cfg GitConfig
}

func NewGit(cfg GitConfig) *GitSegment { return &GitSegment{cfg: cfg} }

func (g *GitSegment) Name() string { return "git" }

func (g *GitSegment) Render() (string, int, int, error) {
	root, err := findGitRoot()
	if err != nil {
		return "", 0, 0, err
	}

	branch, err := readBranch(root)
	if err != nil {
		return "", 0, 0, err
	}

	dirty, err := isDirty()
	if err != nil {
		dirty = false
	}

	text := branch
	if dirty {
		text += " *"
	}

	fg := color.Black
	if g.cfg.FG != nil {
		fg = *g.cfg.FG
	}
	bg := color.Yellow
	if g.cfg.BG != nil {
		bg = *g.cfg.BG
	}

	return text, fg, bg, nil
}

// findGitRoot walks up from cwd to find the nearest .git directory.
func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("not a git repository")
}

// readBranch reads the branch name from .git/HEAD.
func readBranch(root string) (string, error) {
	headPath := filepath.Join(root, ".git", "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("read HEAD: %w", err)
	}
	line := strings.TrimSpace(string(data))
	const prefix = "ref: refs/heads/"
	if strings.HasPrefix(line, prefix) {
		return strings.TrimPrefix(line, prefix), nil
	}
	if len(line) >= 7 {
		return line[:7], nil
	}
	return line, nil
}

// isDirty runs `git status --porcelain` to check for uncommitted changes.
func isDirty() (bool, error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) != "", nil
}
