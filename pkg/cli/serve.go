package cli

import (
	"net/http"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/opac"
	"github.com/secmon-lab/nydus/pkg/adapter"
	"github.com/secmon-lab/nydus/pkg/cli/config"
	"github.com/secmon-lab/nydus/pkg/controller/server"
	"github.com/secmon-lab/nydus/pkg/domain/context/logging"
	"github.com/secmon-lab/nydus/pkg/usecase"
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
			EnvVars:     []string{"NYDUS_ADDR"},
			Usage:       "Address to listen",
			Value:       "127.0.0.1:8080",
			Destination: &addr,
		},
		&cli.StringFlag{
			Name:        "policy-dir",
			Aliases:     []string{"p"},
			EnvVars:     []string{"NYDUS_POLICY_DIR"},
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
		Usage:   "Start nydus server",
		Flags:   flags,
		Action: func(ctx *cli.Context) error {

			logger := logging.Default()
			logger.Info("start nydus server",
				"addr", addr,
				"policyDir", policyDir,
				"azure", azureCfg,
				"gcs", gcsCfg,
			)

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

			mux := server.New(uc)

			logging.Default().Info("starting server", "addr", addr, "policyDir", policyDir)

			httpServer := &http.Server{
				ReadTimeout:  20 * time.Second,
				WriteTimeout: 20 * time.Second,
				IdleTimeout:  120 * time.Second,
				Handler:      mux,
				Addr:         addr,
			}

			if err := httpServer.ListenAndServe(); err != nil {
				return goerr.Wrap(err, "fail to start server")
			}

			return nil
		},
	}
}
