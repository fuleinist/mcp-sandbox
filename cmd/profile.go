package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fuleinist/mcp-sandbox/pkg/profile"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var profileCmd = &cli.Command{
	Name:  "profile",
	Usage: "Manage sandbox profiles",
	Subcommands: []*cli.Command{
		{
			Name:  "list",
			Usage: "List available sandbox profiles",
			Action: func(c *cli.Context) error {
				profiles, err := profile.List()
				if err != nil {
					return fmt.Errorf("failed to list profiles: %w", err)
				}
				if len(profiles) == 0 {
					fmt.Println("No profiles found.")
					return nil
				}
				fmt.Println("Available profiles:")
				for _, p := range profiles {
					fmt.Printf("  %-20s %s\n", p.Name, p.Description)
				}
				return nil
			},
		},
		{
			Name:  "show",
			Usage: "Show details of a specific profile",
			Action: func(c *cli.Context) error {
				if c.Args().Len() == 0 {
					return fmt.Errorf("profile name is required")
				}
				name := c.Args().First()
				p, err := profile.Load(name)
				if err != nil {
					return fmt.Errorf("failed to load profile %q: %w", name, err)
				}
				fmt.Printf("Name:        %s\n", p.Name)
				fmt.Printf("Description: %s\n", p.Description)
				fmt.Printf("Image:       %s\n", p.Image)
				fmt.Printf("Memory:      %s\n", p.Memory)
				fmt.Printf("CPU:         %s\n", p.CPU)
				fmt.Printf("Network:     %v\n", p.AllowNetwork)
				fmt.Printf("Transport:   %s\n", p.Transport)
				if len(p.AllowRead) > 0 {
					fmt.Println("Allow Read:")
					for _, r := range p.AllowRead {
						fmt.Printf("  - %s\n", r)
					}
				}
				if len(p.Env) > 0 {
					fmt.Println("Env:")
					for _, e := range p.Env {
						fmt.Printf("  - %s\n", e)
					}
				}
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "Create a new sandbox profile",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "name", Usage: "Profile name", Required: true},
				&cli.StringFlag{Name: "description", Usage: "Profile description"},
				&cli.StringFlag{Name: "image", Usage: "Docker image", Required: true},
				&cli.StringFlag{Name: "memory", Usage: "Memory limit", Value: "512m"},
				&cli.StringFlag{Name: "cpu", Usage: "CPU limit", Value: "1.0"},
				&cli.BoolFlag{Name: "allow-network", Usage: "Allow network access"},
				&cli.StringFlag{Name: "transport", Usage: "Transport protocol", Value: "stdio"},
				&cli.StringSliceFlag{Name: "allow-read", Usage: "Read-only mount paths"},
				&cli.StringSliceFlag{Name: "deny-write", Usage: "Deny-write paths"},
				&cli.StringSliceFlag{Name: "env", Usage: "Environment variables"},
			},
			Action: func(c *cli.Context) error {
				p := profile.Profile{
					Name:         c.String("name"),
					Description:  c.String("description"),
					Image:        c.String("image"),
					Memory:       c.String("memory"),
					CPU:          c.String("cpu"),
					AllowNetwork: c.Bool("allow-network"),
					Transport:    c.String("transport"),
					AllowRead:    c.StringSlice("allow-read"),
					DenyWrite:    c.StringSlice("deny-write"),
					Env:          c.StringSlice("env"),
				}

				dir, err := profile.ConfigDir()
				if err != nil {
					return err
				}
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create profiles directory: %w", err)
				}

				path := filepath.Join(dir, p.Name+".yaml")
				data, err := yaml.Marshal(&p)
				if err != nil {
					return fmt.Errorf("failed to marshal profile: %w", err)
				}
				if err := os.WriteFile(path, data, 0644); err != nil {
					return fmt.Errorf("failed to write profile: %w", err)
				}

				fmt.Printf("Profile %q created at %s\n", p.Name, path)
				return nil
			},
		},
	},
}
