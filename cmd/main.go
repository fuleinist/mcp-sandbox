package main

import (
	"fmt"
	"os"

	"github.com/fuleinist/mcp-sandbox/pkg/docker"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "mcp-sandbox",
		Usage:       "Safely run MCP servers in ephemeral Docker sandboxes",
		Version:     "0.1.0",
		Description: "Wrap any MCP server command in an ephemeral Docker container with configurable mount points, network restrictions, and resource limits.",
		Commands: []*cli.Command{
			runCmd,
			profileCmd,
		},
		Before: func(c *cli.Context) error {
			if !docker.IsAvailable() {
				return fmt.Errorf("docker is not available — make sure Docker is installed and the daemon is running")
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output in JSON format",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
