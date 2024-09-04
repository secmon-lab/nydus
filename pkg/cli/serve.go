package cli

import "github.com/urfave/cli/v2"

func cmdServe() *cli.Command {
	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Start locust server",
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}
