# mcp-sandbox 🛡️

Docker-based sandbox for safely running untrusted MCP servers with granular filesystem/network permissions.

## Why

MCP servers run as subprocesses inheriting user permissions — a malicious or buggy server can read SSH keys, delete files, exfiltrate env vars. `mcp-sandbox` wraps any MCP server command in an ephemeral Docker container with strict defaults.

## Quick Start

```bash
# Run a Node.js filesystem MCP server sandboxed
mcp-sandbox run \
  --image node:22-alpine \
  --cmd 'npx @modelcontextprotocol/server-filesystem /project' \
  --allow-read /project

# Use a profile
mcp-sandbox run --profile node-fs --cmd 'npx @modelcontextprotocol/server-filesystem /project'

# SSE transport
mcp-sandbox run \
  --image python:3.12 \
  --cmd 'python -m mcp_server_postgres' \
  --transport sse --port 3100 \
  --allow-network
```

## Installation

```bash
go install github.com/fuleinist/mcp-sandbox/cmd/mcp-sandbox@latest
```

Or download a binary from the [releases page](https://github.com/fuleinist/mcp-sandbox/releases).

## Usage

### Command: `run`

```
mcp-sandbox run [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--image, -i` | (required) | Docker image |
| `--cmd, -c` | (required) | Command to run |
| `--allow-read` | | Read-only mount paths |
| `--allow-network` | `false` | Enable network access |
| `--memory` | `512m` | Memory limit |
| `--cpu` | `1.0` | CPU limit |
| `--transport` | `stdio` | `stdio` or `sse` |
| `--port` | | Host port (SSE only) |
| `--env` | | Environment variables |
| `--profile, -p` | | Sandbox profile name |
| `--dry-run` | | Print docker command |
| `--verbose, -v` | | Show docker commands |

### Command: `profile`

```
mcp-sandbox profile list     # List available profiles
mcp-sandbox profile show <n> # Show profile details
```

## Profiles

Save reusable configs to `~/.config/mcp-sandbox/profiles/`:

```yaml
# ~/.config/mcp-sandbox/profiles/my-app.yaml
name: my-app
description: My custom MCP server
image: node:22-alpine
memory: 1g
cpu: 2.0
allow_network: true
transport: stdio
allow_read:
  - /project
  - /data
env:
  - NODE_ENV=production
  - LOG_LEVEL=debug
```

Built-in profiles: `node-fs`, `python-default`, `deno-basic`, `postgres-sse`.

## Security

- ❌ No network access by default
- ❌ Read-only root filesystem
- ❌ No new privileges (`no-new-privileges`)
- ❌ All Linux capabilities dropped (`--cap-drop ALL`)
- 👤 Runs as non-root user (UID 1000)
- ✅ Explicit read/write mount opt-in

## License

MIT
