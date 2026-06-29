package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsAvailable checks if Docker is running by calling `docker info`.
func IsAvailable() bool {
	cmd := exec.Command("docker", "info", "--format", "{{.ServerVersion}}")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}

// BuildRunArgs constructs the docker run argument list from the given parameters.
func BuildRunArgs(
	image string,
	cmd string,
	allowRead []string,
	allowNet bool,
	memory string,
	cpu string,
	transport string,
	port int,
	env []string,
	autoRemove bool,
) []string {
	args := []string{"run"}

	if autoRemove {
		args = append(args, "--rm")
	}

	// Run as non-root
	args = append(args, "--user", "1000:1000")

	// Resource limits
	if memory != "" {
		args = append(args, "--memory", memory)
	}
	if cpu != "" {
		args = append(args, "--cpus", cpu)
	}

	// Network
	if !allowNet {
		args = append(args, "--network", "none")
	}

	// Read-only root filesystem
	args = append(args, "--read-only")

	// Read-only mounts
	for _, p := range allowRead {
		args = append(args, "--mount", fmt.Sprintf("type=bind,source=%s,target=%s,readonly", p, p))
	}

	// Environment variables
	for _, e := range env {
		args = append(args, "-e", e)
	}

	// Port mapping (SSE)
	if transport == "sse" && port > 0 {
		args = append(args, "-p", fmt.Sprintf("%d:%d", port, port))
	}

	// Security options
	args = append(args, "--security-opt", "no-new-privileges:true")
	args = append(args, "--cap-drop", "ALL")

	// Image
	args = append(args, image)

	// Command (split by spaces for simplicity; users can use shell wrapper for complex commands)
	args = append(args, "sh", "-c", cmd)

	return args
}

// PrintRunCommand prints the docker run command for dry-run/verbose mode.
func PrintRunCommand(args []string) {
	fmt.Printf("docker %s\n", strings.Join(args, " "))
}


