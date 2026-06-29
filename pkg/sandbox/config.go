package sandbox

import "github.com/fuleinist/mcp-sandbox/pkg/profile"

// Config holds all configuration for a sandboxed MCP server run.
type Config struct {
	Image      string
	Cmd        string
	AllowRead  []string
	DenyWrite  []string
	AllowNet   bool
	DenyNet    bool
	Memory     string
	CPU        string
	Transport  string // "stdio" or "sse"
	Port       int
	Env        []string
	AutoRemove bool
	DryRun     bool
	Verbose    bool
}

// ApplyProfile overlays profile values onto the config.
// CLI flags take precedence over profile values.
func (c *Config) ApplyProfile(p *profile.Profile) {
	if c.Image == "" && p.Image != "" {
		c.Image = p.Image
	}
	if c.Memory == "" && p.Memory != "" {
		c.Memory = p.Memory
	}
	if c.CPU == "" && p.CPU != "" {
		c.CPU = p.CPU
	}
	if !c.AllowNet && p.AllowNetwork {
		c.AllowNet = p.AllowNetwork
	}
	if c.Transport == "" && p.Transport != "" {
		c.Transport = p.Transport
	}
	// Merge read paths
	if len(p.AllowRead) > 0 {
		existing := make(map[string]bool)
		for _, r := range c.AllowRead {
			existing[r] = true
		}
		for _, r := range p.AllowRead {
			if !existing[r] {
				c.AllowRead = append(c.AllowRead, r)
			}
		}
	}
	// Merge deny-write paths
	if len(p.DenyWrite) > 0 {
		existing := make(map[string]bool)
		for _, r := range c.DenyWrite {
			existing[r] = true
		}
		for _, r := range p.DenyWrite {
			if !existing[r] {
				c.DenyWrite = append(c.DenyWrite, r)
			}
		}
	}
	// Merge env vars
	if len(p.Env) > 0 {
		existing := make(map[string]bool)
		for _, e := range c.Env {
			existing[e] = true
		}
		for _, e := range p.Env {
			if !existing[e] {
				c.Env = append(c.Env, e)
			}
		}
	}
}
