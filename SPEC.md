# mcp-sandbox — SPEC v1

## Mission
A CLI tool that wraps any MCP server command in an ephemeral Docker container with configurable mount points, network restrictions, and resource limits. Supports both stdio and SSE transport.

## Why
MCP servers run as subprocesses inheriting user permissions — a malicious or buggy server can read SSH keys, delete files, exfiltrate env vars. There's no safety layer.

## Usage
```bash
# Basic — run a Node.js MCP server in a sandbox
mcp-sandbox --image node:22 --cmd 'npx @modelcontextprotocol/server-filesystem' --allow-read /project

# With granular permissions
mcp-sandbox --image python:3.12 --cmd 'python -m mcp_server_postgres' \
  --allow-read /data \
  --deny-write \
  --deny-network \
  --memory 512m \
  --cpu 1.0

# With a sandbox profile (YAML)
mcp-sandbox --profile node-fs

# SSE transport
mcp-sandbox --image node:22 --cmd 'node server.mjs' --transport sse --port 3100
```

## Acceptance Criteria

### CLI (MVP)
- [ ] `mcp-sandbox` binary with subcommands: `run`, `profile`, `list`, `help`
- [ ] `--image` flag: Docker image to use (required)
- [ ] `--cmd` flag: command to run inside container (required)
- [ ] `--allow-read` / `--deny-write`: filesystem mount rules (repeatable)
- [ ] `--allow-network` / `--deny-network`: network access toggle (default: allow)
- [ ] `--memory`: memory limit (e.g. `512m`, `2g`)
- [ ] `--cpu`: CPU limit (e.g. `1.0`, `0.5`)
- [ ] `--transport`: `stdio` (default) or `sse`
- [ ] `--port`: host port mapping (for SSE transport)
- [ ] `--profile`: load sandbox profile from YAML file
- [ ] `--env`: pass environment variables into container (repeatable)
- [ ] `--rm`: auto-remove container on exit (default: true)
- [ ] Stdio mode: pipes stdin/stdout/stderr between host and container
- [ ] SSE mode: exposes container port on host, forwards traffic
- [ ] Exit code passthrough from container

### Profiles
- [ ] YAML config file format (default: `~/.config/mcp-sandbox/profiles/`)
- [ ] Built-in profiles: `node-fs`, `python-default`, `deno-basic`
- [ ] `mcp-sandbox profile list` — list available profiles
- [ ] `mcp-sandbox profile show <name>` — show profile details

### Security
- [ ] Default: no network, no write access, read-only root
- [ ] Container runs as non-root user by default
- [ ] Read mounts are read-only bind mounts
- [ ] Write mounts are explicit opt-in
- [ ] Network off by default, opt-in with `--allow-network`
- [ ] Resource limits prevent fork bombs / memory DoS

### Developer Experience
- [ ] Clear error messages (Docker not running, image not found, etc.)
- [ ] JSON output mode (`--json`) for programmatic use
- [ ] Dry-run mode (`--dry-run`) prints docker run command without executing
- [ ] Verbose mode (`-v`) shows underlying docker commands

## Tech Stack
- **Language:** Go (single binary, good Docker SDK support)
- **Docker SDK:** `github.com/docker/docker/client` (official Go SDK)
- **CLI framework:** `github.com/urfave/cli/v2`
- **Config:** YAML via `gopkg.in/yaml.v3`
- **Build:** Go 1.22+, `go build` for single binary

## Architecture

```
mcp-sandbox
├── cmd/                  # CLI entry points
│   ├── run.go            # `mcp-sandbox run` subcommand
│   ├── profile.go        # `mcp-sandbox profile` subcommand
│   └── main.go           # root CLI setup
├── pkg/
│   ├── sandbox/          # Core sandbox logic
│   │   ├── container.go  # Docker container lifecycle
│   │   ├── stdio.go      # Stdio transport handler
│   │   ├── sse.go        # SSE transport handler
│   │   └── config.go     # Sandbox configuration
│   ├── profile/          # Profile management
│   │   ├── loader.go     # YAML profile loading
│   │   └── builtin.go    # Built-in profiles
│   └── docker/           # Docker client wrapper
│       └── client.go     # Docker SDK abstraction
├── profiles/             # Built-in profile YAML files
├── SPEC.md               # This file
└── README.md             # User-facing docs
```

## Out of Scope (v1)
- Windows container support (Linux containers only for MVP)
- Kubernetes / podman support (Docker-only for MVP)
- MCP protocol inspection / validation
- Web UI / dashboard
- Remote Docker daemon over SSH
