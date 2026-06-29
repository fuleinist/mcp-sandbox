package sandbox

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fuleinist/mcp-sandbox/pkg/docker"
)

// Runner manages the lifecycle of a sandboxed MCP server.
type Runner struct {
	cfg Config
}

// RunResult represents the outcome of a sandbox run.
type RunResult struct {
	ExitCode    int    `json:"exit_code"`
	ContainerID string `json:"container_id,omitempty"`
	Transport   string `json:"transport"`
	Port        int    `json:"port,omitempty"`
}

// NewRunner creates a new sandbox runner.
func NewRunner(cfg Config) *Runner {
	return &Runner{cfg: cfg}
}

// Run executes the sandboxed MCP server.
// For stdio transport, it pipes stdin/stdout/stderr between host and container.
// For SSE transport, it runs the container in detached mode with port forwarding.
func (r *Runner) Run(ctx context.Context) (int, error) {
	args := docker.BuildRunArgs(
		r.cfg.Image,
		r.cfg.Cmd,
		r.cfg.AllowRead,
		r.cfg.AllowNet,
		r.cfg.Memory,
		r.cfg.CPU,
		r.cfg.Transport,
		r.cfg.Port,
		r.cfg.Env,
		r.cfg.AutoRemove,
	)

	if r.cfg.DryRun {
		docker.PrintRunCommand(args)
		return 0, nil
	}

	if r.cfg.Verbose {
		docker.PrintRunCommand(args)
	}

	switch r.cfg.Transport {
	case "stdio":
		return r.runStdio(ctx, args)
	case "sse":
		return r.runSSE(ctx, args)
	default:
		return 1, fmt.Errorf("unsupported transport: %s", r.cfg.Transport)
	}
}

// runStdio runs the container with stdin/stdout/stderr connected directly.
func (r *Runner) runStdio(ctx context.Context, args []string) (int, error) {
	cmd := exec.CommandContext(ctx, "docker", args...)

	// Connect stdio
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}
		return 1, fmt.Errorf("container run failed: %w", err)
	}

	return cmd.ProcessState.ExitCode(), nil
}

// runSSE runs the container in detached mode with port forwarding.
func (r *Runner) runSSE(ctx context.Context, args []string) (int, error) {
	// For SSE, we run detached and stream logs
	detachArgs := make([]string, len(args))
	copy(detachArgs, args)
	detachArgs = append(detachArgs[:1], append([]string{"-d"}, detachArgs[1:]...)...)

	cmd := exec.CommandContext(ctx, "docker", detachArgs...)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return 1, fmt.Errorf("failed to start container: %w", err)
	}

	id := strings.TrimSpace(string(out))
	fmt.Printf("Container started: %s\n", id)
	fmt.Printf("MCP server listening on port %d (SSE)\n", r.cfg.Port)
	fmt.Println("Press Ctrl+C to stop.")

	// Stream logs until interrupted
	logCmd := exec.CommandContext(ctx, "docker", "logs", "-f", id)
	logCmd.Stdout = os.Stdout
	logCmd.Stderr = os.Stderr
	_ = logCmd.Run()

	// Cleanup
	rmCmd := exec.CommandContext(ctx, "docker", "rm", "-f", id)
	_ = rmCmd.Run()

	return 0, nil
}

// OutputJSON prints the run result as JSON.
func OutputJSON(result RunResult) {
	data, _ := json.Marshal(result)
	fmt.Println(string(data))
}


