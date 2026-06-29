package profile

// Profile defines a reusable sandbox configuration.
type Profile struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Image        string   `yaml:"image"`
	Memory       string   `yaml:"memory"`
	CPU          string   `yaml:"cpu"`
	AllowNetwork bool     `yaml:"allow_network"`
	Transport    string   `yaml:"transport"`
	AllowRead    []string `yaml:"allow_read"`
	DenyWrite    []string `yaml:"deny_write"`
	Env          []string `yaml:"env"`
}
