package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigDir returns the default config directory for profiles.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "mcp-sandbox", "profiles"), nil
}

// Load loads a profile by name. Checks built-in profiles first, then user profiles.
func Load(name string) (*Profile, error) {
	// Check built-in first
	if p, ok := builtins[name]; ok {
		return &p, nil
	}

	// Check user profiles directory
	dir, err := ConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, name+".yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Also try .yml
		path = filepath.Join(dir, name+".yml")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, fmt.Errorf("profile %q not found (checked built-ins and %s)", name, dir)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile %q: %w", name, err)
	}

	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse profile %q: %w", name, err)
	}

	if p.Name == "" {
		p.Name = name
	}

	return &p, nil
}

// List returns all available profiles (built-in + user).
func List() ([]Profile, error) {
	var profiles []Profile

	// Add built-in profiles
	for _, p := range builtins {
		profiles = append(profiles, p)
	}

	// Scan user profiles directory
	dir, err := ConfigDir()
	if err != nil {
		return profiles, nil // return built-ins even if config dir fails
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return profiles, nil // no user profiles, return built-ins
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var p Profile
		if err := yaml.Unmarshal(data, &p); err != nil {
			continue
		}
		if p.Name == "" {
			p.Name = entry.Name()[:len(entry.Name())-len(ext)]
		}
		profiles = append(profiles, p)
	}

	return profiles, nil
}
