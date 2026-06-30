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

### CLI (MVP) ✅
- [x] `mcp-sandbox` binary with subcommands: `run`, `profile`, `help`
- [x] `--image` flag: Docker image to use (required)
- [x] `--cmd` flag: command to run inside container (required)
- [x] `--allow-read` / `--deny-write`: filesystem mount rules (repeatable)
- [x] `--allow-network` / `--deny-network`: network access toggle (default: deny)
- [x] `--memory`: memory limit (e.g. `512m`, `2g`)
- [x] `--cpu`: CPU limit (e.g. `1.0`, `0.5`)
- [x] `--transport`: `stdio` (default) or `sse`
- [x] `--port`: host port mapping (for SSE transport)
- [x] `--profile`: load sandbox profile from YAML file or built-in name
- [x] `--env`: pass environment variables into container (repeatable)
- [x] `--rm`: auto-remove container on exit (default: true)
- [x] Stdio mode: pipes stdin/stdout/stderr between host and container
- [x] SSE mode: exposes container port on host, forwards traffic
- [x] Exit code passthrough from container
- [x] `--dry-run`: print docker command without executing
- [x] `--verbose` / `-v`: show underlying docker commands
- [x] `--json`: JSON output mode for programmatic use

### Profiles ✅
- [x] YAML config file format (default: `~/.config/mcp-sandbox/profiles/`)
- [x] Built-in profiles: `node-fs`, `python-default`, `deno-basic`, `postgres-sse`
- [x] `mcp-sandbox profile list` — list available profiles
- [x] `mcp-sandbox profile show <name>` — show profile details
- [x] `mcp-sandbox profile create` — create a new profile

### Security ✅
- [x] Default: no network, no write access, read-only root
- [x] Container runs as non-root user by default (UID 1000)
- [x] Read mounts are read-only bind mounts
- [x] Write mounts are explicit opt-in
- [x] Network off by default, opt-in with `--allow-network`
- [x] Resource limits (`--memory`, `--cpus`)
- [x] All Linux capabilities dropped (`--cap-drop ALL`)
- [x] No new privileges (`--security-opt no-new-privileges:true`)

### Developer Experience ✅
- [x] Clear error messages (Docker not running, image not found, etc.)
- [x] JSON output mode (`--json`) for programmatic use
- [x] Dry-run mode (`--dry-run`) prints docker run command without executing
- [x] Verbose mode (`-v`) shows underlying docker commands
- [x] `--json` flag on `run` subcommand
- [x] 11 unit tests covering all core functionality

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
