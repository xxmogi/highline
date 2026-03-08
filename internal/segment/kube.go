package segment

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/highline/internal/color"
)

// KubeConfig holds optional color configuration for KubeSegment.
type KubeConfig struct {
	FG *int // nil = color.White
	BG *int // nil = color.Cyan
}

// KubeSegment renders the current Kubernetes context.
type KubeSegment struct {
	cfg KubeConfig
}

func NewKube(cfg KubeConfig) *KubeSegment { return &KubeSegment{cfg: cfg} }

func (k *KubeSegment) Name() string { return "kube" }

func (k *KubeSegment) Render() (string, int, int, error) {
	ctx, err := currentKubeContext()
	if err != nil {
		return "", 0, 0, err
	}

	fg := color.White
	if k.cfg.FG != nil {
		fg = *k.cfg.FG
	}
	bg := color.Cyan
	if k.cfg.BG != nil {
		bg = *k.cfg.BG
	}

	return "\u2388 " + ctx, fg, bg, nil
}

// currentKubeContext reads the current-context field from the kubeconfig file.
func currentKubeContext() (string, error) {
	path := kubeConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("kubeconfig not found: %w", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "current-context:") {
			ctx := strings.TrimSpace(strings.TrimPrefix(trimmed, "current-context:"))
			ctx = strings.Trim(ctx, `"'`)
			if ctx == "" || ctx == "null" {
				return "", errors.New("no current context set")
			}
			return ctx, nil
		}
	}
	return "", errors.New("current-context not found in kubeconfig")
}

// kubeConfigPath returns the path to the kubeconfig file.
func kubeConfigPath() string {
	if v := os.Getenv("KUBECONFIG"); v != "" {
		parts := strings.SplitN(v, ":", 2)
		return parts[0]
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}
