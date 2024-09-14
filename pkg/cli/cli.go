package cli

import (
	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/nydus/pkg/cli/config"
	"github.com/secmon-as-code/nydus/pkg/domain/context/logging"
	"github.com/urfave/cli/v2"
)

func Run(argv []string) error {
	var flags []cli.Flag

	var loggingCfg config.Logging
	flags = append(flags, loggingCfg.Flags()...)

	app := cli.App{
		Name:  "nydus",
		Flags: flags,
		Commands: []*cli.Command{
			cmdServe(),
		},
		Before: func(ctx *cli.Context) error {
			logger, err := loggingCfg.NewLogger()
			if err != nil {
				return goerr.Wrap(err, "fail to create logger")
			}

			logging.SetDefault(logger)
			return nil
		},
	}

	if err := app.Run(argv); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
