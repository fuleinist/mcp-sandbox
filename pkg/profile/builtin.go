package profile

// builtins contains the built-in sandbox profiles.
var builtins = map[string]Profile{
	"node-fs": {
		Name:        "node-fs",
		Description: "Node.js MCP server with filesystem access",
		Image:       "node:22-alpine",
		Memory:      "512m",
		CPU:         "1.0",
		AllowNetwork: false,
		Transport:   "stdio",
		AllowRead:   []string{"/project"},
		DenyWrite:   []string{"/root/.ssh", "/root/.gitconfig"},
		Env:         []string{"NODE_ENV=production"},
	},
	"python-default": {
		Name:        "python-default",
		Description: "Python MCP server with standard libraries",
		Image:       "python:3.12-alpine",
		Memory:      "512m",
		CPU:         "1.0",
		AllowNetwork: false,
		Transport:   "stdio",
		AllowRead:   []string{"/project"},
	},
	"deno-basic": {
		Name:        "deno-basic",
		Description: "Deno MCP server with network access (Deno needs net for module fetching)",
		Image:       "denoland/deno:alpine",
		Memory:      "512m",
		CPU:         "1.0",
		AllowNetwork: true,
		Transport:   "stdio",
		AllowRead:   []string{"/project"},
	},
	"postgres-sse": {
		Name:        "postgres-sse",
		Description: "Postgres MCP server exposed via SSE",
		Image:       "python:3.12-alpine",
		Memory:      "1g",
		CPU:         "1.0",
		AllowNetwork: true,
		Transport:   "sse",
		AllowRead:   []string{"/data"},
	},
}
