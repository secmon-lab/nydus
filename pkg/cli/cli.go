package cli

import (
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func Run(argv []string) error {
	app := cli.App{
		Name:     "locust",
		Commands: []*cli.Command{},
	}

	if err := app.Run(argv); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
