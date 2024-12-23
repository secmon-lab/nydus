package config

import (
	"log/slog"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/nydus/pkg/adapter/abs"
	"github.com/secmon-lab/nydus/pkg/domain/context/logging"
	"github.com/urfave/cli/v2"
)

type Azure struct {
	enable       bool
	tenantID     string
	clientID     string
	clientSecret string
}

func (x *Azure) Flags() []cli.Flag {
	const category = "Azure Blob Storage"

	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "enable-azure",
			Usage:       "Enable Azure Blob Storage",
			Category:    category,
			EnvVars:     []string{"NYDUS_ENABLE_AZURE"},
			Destination: &x.enable,
		},

		&cli.StringFlag{
			Name:        "azure-tenant-id",
			Usage:       "Azure tenant ID",
			Category:    category,
			EnvVars:     []string{"NYDUS_AZURE_TENANT_ID"},
			Destination: &x.tenantID,
		},
		&cli.StringFlag{
			Name:        "azure-client-id",
			Usage:       "Azure client ID",
			Category:    category,
			EnvVars:     []string{"NYDUS_AZURE_CLIENT_ID"},
			Destination: &x.clientID,
		},
		&cli.StringFlag{
			Name:        "azure-client-secret",
			Usage:       "Azure client secret",
			Category:    category,
			EnvVars:     []string{"NYDUS_AZURE_CLIENT_SECRET"},
			Destination: &x.clientSecret,
		},
	}
}

func (x Azure) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Bool("enable", x.enable),
		slog.String("tenantID", x.tenantID),
		slog.String("clientID", x.clientID),
		slog.Int("clientSecret(len)", len(x.clientSecret)),
	)
}

func (x *Azure) NewClient() (*abs.Client, error) {
	if !x.enable {
		if x.tenantID != "" || x.clientID != "" || x.clientSecret != "" {
			logging.Default().Warn("Azure configuration is ignored because Azure is disabled")
		}
		return nil, nil
	}

	if x.tenantID == "" {
		return nil, goerr.New("Azure tenant ID is required")
	}
	if x.clientID == "" {
		return nil, goerr.New("Azure client ID is required")
	}
	if x.clientSecret == "" {
		return nil, goerr.New("Azure client secret is required")
	}

	return abs.New(x.clientID, x.clientSecret, x.tenantID)
}
