package sandbox

import (
	"testing"

	"github.com/fuleinist/mcp-sandbox/pkg/profile"
)

func TestApplyProfile_EmptyConfig(t *testing.T) {
	cfg := Config{}
	p := &profile.Profile{
		Name:    "test",
		Image:   "node:22",
		Memory:  "1g",
		CPU:     "2.0",
		AllowNetwork: true,
		Transport: "sse",
		AllowRead: []string{"/data"},
		Env:      []string{"FOO=bar"},
	}
	cfg.ApplyProfile(p)

	if cfg.Image != "node:22" {
		t.Errorf("expected image node:22, got %s", cfg.Image)
	}
	if cfg.Memory != "1g" {
		t.Errorf("expected memory 1g, got %s", cfg.Memory)
	}
	if cfg.CPU != "2.0" {
		t.Errorf("expected cpu 2.0, got %s", cfg.CPU)
	}
	if !cfg.AllowNet {
		t.Error("expected AllowNet true")
	}
	if cfg.Transport != "sse" {
		t.Errorf("expected transport sse, got %s", cfg.Transport)
	}
	if len(cfg.AllowRead) != 1 || cfg.AllowRead[0] != "/data" {
		t.Errorf("expected allow-read [/data], got %v", cfg.AllowRead)
	}
	if len(cfg.Env) != 1 || cfg.Env[0] != "FOO=bar" {
		t.Errorf("expected env [FOO=bar], got %v", cfg.Env)
	}
}

func TestApplyProfile_CLIPrecedence(t *testing.T) {
	cfg := Config{
		Image: "python:3.12",
		CPU:   "4.0",
	}
	p := &profile.Profile{
		Image: "node:22",
		CPU:   "2.0",
	}
	cfg.ApplyProfile(p)

	// CLI values should not be overwritten
	if cfg.Image != "python:3.12" {
		t.Errorf("expected image python:3.12 (CLI precedence), got %s", cfg.Image)
	}
	if cfg.CPU != "4.0" {
		t.Errorf("expected cpu 4.0 (CLI precedence), got %s", cfg.CPU)
	}
}

func TestApplyProfile_MergeDenyWrite(t *testing.T) {
	cfg := Config{
		DenyWrite: []string{"/etc"},
	}
	p := &profile.Profile{
		DenyWrite: []string{"/root/.ssh", "/etc"},
	}
	cfg.ApplyProfile(p)

	if len(cfg.DenyWrite) != 2 {
		t.Errorf("expected 2 unique deny-write paths, got %d: %v", len(cfg.DenyWrite), cfg.DenyWrite)
	}
}

func TestApplyProfile_MergeReadPaths(t *testing.T) {
	cfg := Config{
		AllowRead: []string{"/project"},
	}
	p := &profile.Profile{
		AllowRead: []string{"/data", "/project"},
	}
	cfg.ApplyProfile(p)

	if len(cfg.AllowRead) != 2 {
		t.Errorf("expected 2 unique read paths, got %d: %v", len(cfg.AllowRead), cfg.AllowRead)
	}
}

func TestApplyProfile_MergeEnv(t *testing.T) {
	cfg := Config{
		Env: []string{"A=1"},
	}
	p := &profile.Profile{
		Env: []string{"B=2", "A=1"},
	}
	cfg.ApplyProfile(p)

	if len(cfg.Env) != 2 {
		t.Errorf("expected 2 unique env vars, got %d: %v", len(cfg.Env), cfg.Env)
	}
}
