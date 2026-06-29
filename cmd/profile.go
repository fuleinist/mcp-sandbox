package main

import (
	"fmt"

	"github.com/fuleinist/mcp-sandbox/pkg/profile"
	"github.com/urfave/cli/v2"
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
	},
}
