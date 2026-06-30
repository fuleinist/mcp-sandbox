package docker

import (
	"strings"
	"testing"
)

func TestBuildRunArgs_Basic(t *testing.T) {
	args := BuildRunArgs(
		"node:22-alpine",
		"npx @modelcontextprotocol/server-filesystem /project",
		[]string{"/project"},
		nil,
		false,
		false,
		"512m",
		"1.0",
		"stdio",
		0,
		nil,
		true,
	)

	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "node:22-alpine") {
		t.Error("expected image in args")
	}
	if !strings.Contains(joined, "--rm") {
		t.Error("expected --rm flag")
	}
	if !strings.Contains(joined, "--read-only") {
		t.Error("expected --read-only flag")
	}
	if !strings.Contains(joined, "--network none") {
		t.Error("expected --network none for no network")
	}
	if !strings.Contains(joined, "--user 1000:1000") {
		t.Error("expected non-root user")
	}
	if !strings.Contains(joined, "--cap-drop ALL") {
		t.Error("expected --cap-drop ALL")
	}
	if !strings.Contains(joined, "--security-opt no-new-privileges:true") {
		t.Error("expected no-new-privileges")
	}
	if !strings.Contains(joined, "--memory 512m") {
		t.Error("expected memory limit")
	}
	if !strings.Contains(joined, "--cpus 1.0") {
		t.Error("expected cpu limit")
	}
}

func TestBuildRunArgs_SSEWithPort(t *testing.T) {
	args := BuildRunArgs(
		"python:3.12",
		"python -m mcp_server",
		[]string{"/data"},
		nil,
		true,
		false,
		"1g",
		"2.0",
		"sse",
		3100,
		[]string{"PORT=3100"},
		true,
	)

	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-p 3100:3100") {
		t.Error("expected port mapping for SSE")
	}
	if strings.Contains(joined, "--network none") {
		t.Error("should not have --network none when allowNet is true")
	}
	if !strings.Contains(joined, "-e PORT=3100") {
		t.Error("expected env var")
	}
}

func TestBuildRunArgs_ReadOnlyMounts(t *testing.T) {
	args := BuildRunArgs(
		"node:22",
		"node server.js",
		[]string{"/project", "/data/config"},
		nil,
		false,
		false,
		"",
		"",
		"stdio",
		0,
		nil,
		false,
	)

	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "source=/project,target=/project,readonly") {
		t.Error("expected /project read-only mount")
	}
	if !strings.Contains(joined, "source=/data/config,target=/data/config,readonly") {
		t.Error("expected /data/config read-only mount")
	}
}

func TestBuildRunArgs_DenyWrite(t *testing.T) {
	args := BuildRunArgs(
		"node:22",
		"node server.js",
		nil,
		[]string{"/etc", "/root/.ssh"},
		false,
		false,
		"",
		"",
		"stdio",
		0,
		nil,
		true,
	)

	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--tmpfs /etc:noexec,nosuid,size=1m") {
		t.Error("expected tmpfs mount for /etc deny-write")
	}
	if !strings.Contains(joined, "--tmpfs /root/.ssh:noexec,nosuid,size=1m") {
		t.Error("expected tmpfs mount for /root/.ssh deny-write")
	}
}

func TestPrintRunCommand(t *testing.T) {
	// Just ensure it doesn't panic
	PrintRunCommand([]string{"run", "--rm", "node:22", "node", "server.js"})
}

func TestBuildRunArgs_DenyNet(t *testing.T) {
	args := BuildRunArgs(
		"node:22",
		"node server.js",
		nil,
		nil,
		true,
		true,
		"",
		"",
		"stdio",
		0,
		nil,
		true,
	)

	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--network none") {
		t.Error("expected --network none when denyNet is true even if allowNet is true")
	}
}
