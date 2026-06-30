package main

import (
	"fmt"
	"os"

	"github.com/fuleinist/mcp-sandbox/pkg/profile"
	"github.com/fuleinist/mcp-sandbox/pkg/sandbox"
	"github.com/urfave/cli/v2"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Run an MCP server in a sandboxed Docker container",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "image",
			Usage:    "Docker image to use (e.g. `node:22`, `python:3.12`)",
			Required: true,
			Aliases:  []string{"i"},
		},
		&cli.StringFlag{
			Name:     "cmd",
			Usage:    "Command to run inside the container",
			Required: true,
			Aliases:  []string{"c"},
		},
		&cli.StringSliceFlag{
			Name:  "allow-read",
			Usage: "Paths to mount as read-only (repeatable)",
		},
		&cli.StringSliceFlag{
			Name:  "deny-write",
			Usage: "Paths to explicitly deny write access (repeatable)",
		},
		&cli.BoolFlag{
			Name:  "allow-network",
			Usage: "Allow network access (default: false)",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "deny-network",
			Usage: "Explicitly deny network access (overrides --allow-network)",
		},
		&cli.StringFlag{
			Name:  "memory",
			Usage: "Memory limit (e.g. `512m`, `2g`)",
			Value: "512m",
		},
		&cli.StringFlag{
			Name:  "cpu",
			Usage: "CPU limit (e.g. `1.0`, `0.5`)",
			Value: "1.0",
		},
		&cli.StringFlag{
			Name:  "transport",
			Usage: "Transport protocol: `stdio` or `sse`",
			Value: "stdio",
		},
		&cli.IntFlag{
			Name:  "port",
			Usage: "Host port to map (for SSE transport)",
		},
		&cli.StringFlag{
			Name:  "profile",
			Usage: "Load sandbox profile from YAML file or built-in name",
			Aliases: []string{"p"},
		},
		&cli.StringSliceFlag{
			Name:  "env",
			Usage: "Environment variables to pass into container (KEY=VALUE, repeatable)",
		},
		&cli.BoolFlag{
			Name:  "rm",
			Usage: "Auto-remove container on exit",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Print the docker command without executing",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Show underlying docker commands",
		},
		&cli.BoolFlag{
			Name:  "json",
			Usage: "Output in JSON format",
		},
	},
	Action: func(c *cli.Context) error {
		cfg := sandbox.Config{
			Image:       c.String("image"),
			Cmd:         c.String("cmd"),
			AllowRead:   c.StringSlice("allow-read"),
			DenyWrite:   c.StringSlice("deny-write"),
			AllowNet:    c.Bool("allow-network"),
			DenyNet:     c.Bool("deny-network"),
			Memory:      c.String("memory"),
			CPU:         c.String("cpu"),
			Transport:   c.String("transport"),
			Port:        c.Int("port"),
			Env:         c.StringSlice("env"),
			AutoRemove:  c.Bool("rm"),
			DryRun:      c.Bool("dry-run"),
			Verbose:     c.Bool("verbose"),
		}

		// Load profile if specified
		if profileName := c.String("profile"); profileName != "" {
			p, err := profile.Load(profileName)
			if err != nil {
				return fmt.Errorf("failed to load profile %q: %w", profileName, err)
			}
			cfg.ApplyProfile(p)
		}

		// Validate transport
		if cfg.Transport != "stdio" && cfg.Transport != "sse" {
			return fmt.Errorf("invalid transport %q: must be 'stdio' or 'sse'", cfg.Transport)
		}

		// SSE requires a port
		if cfg.Transport == "sse" && cfg.Port == 0 {
			return fmt.Errorf("--port is required when using SSE transport")
		}

		jsonOutput := c.Bool("json")

		runner := sandbox.NewRunner(cfg)
		exitCode, err := runner.Run(c.Context)
		if err != nil {
			if jsonOutput {
				sandbox.OutputJSON(sandbox.RunResult{ExitCode: 1})
			}
			return err
		}

		if jsonOutput {
			sandbox.OutputJSON(sandbox.RunResult{
				ExitCode:  exitCode,
				Transport: cfg.Transport,
				Port:      cfg.Port,
			})
		}

		os.Exit(exitCode)
		return nil
	},
}
