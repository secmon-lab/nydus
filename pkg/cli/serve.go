package cli

import (
	"net/http"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/opac"
	"github.com/secmon-as-code/locust/pkg/adapter"
	"github.com/secmon-as-code/locust/pkg/cli/config"
	"github.com/secmon-as-code/locust/pkg/controller/server"
	"github.com/secmon-as-code/locust/pkg/domain/context/logging"
	"github.com/secmon-as-code/locust/pkg/usecase"
	"github.com/urfave/cli/v2"
)

func cmdServe() *cli.Command {
	var (
		addr      string
		policyDir string
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "addr",
			Aliases:     []string{"a"},
			EnvVars:     []string{"LOCUST_ADDR"},
			Usage:       "Address to listen",
			Value:       "127.0.0.1:8080",
			Destination: &addr,
		},
		&cli.StringFlag{
			Name:        "policy-dir",
			Aliases:     []string{"p"},
			EnvVars:     []string{"LOCUST_POLICY_DIR"},
			Usage:       "Directory path of policy files",
			Value:       "policy",
			Destination: &policyDir,
			Required:    true,
		},
	}

	var azureCfg config.Azure
	flags = append(flags, azureCfg.Flags()...)

	var gcsCfg config.GoogleCloudStorage
	flags = append(flags, gcsCfg.Flags()...)

	return &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Start locust server",
		Flags:   flags,
		Action: func(ctx *cli.Context) error {
			policy, err := opac.New(opac.Files(policyDir))
			if err != nil {
				return goerr.Wrap(err, "fail to load policy files")
			}

			adaptorOptions := []adapter.Option{
				adapter.WithPolicy(policy),
			}

			// Setup Azure Blob Storage client
			if client, err := azureCfg.NewClient(); err != nil {
				return goerr.Wrap(err, "fail to create Azure Blob Storage client")
			} else if client != nil {
				adaptorOptions = append(adaptorOptions, adapter.WithAzureBlobStorage(client))
			}

			// Setup Google Cloud Storage client
			if client, err := gcsCfg.NewClient(); err != nil {
				return goerr.Wrap(err, "fail to create Google Cloud Storage client")
			} else if client != nil {
				adaptorOptions = append(adaptorOptions, adapter.WithGoogleCloudStorage(client))
			}

			clients := adapter.New(adaptorOptions...)

			uc := usecase.New(clients)

			srv := server.New(uc)

			logging.Default().Info("starting server", "addr", addr, "policyDir", policyDir)
			if err := http.ListenAndServe(addr, srv); err != nil {
				return goerr.Wrap(err, "fail to start server")
			}

			return nil
		},
	}
}
