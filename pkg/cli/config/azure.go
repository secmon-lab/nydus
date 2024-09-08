package config

import (
	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/locust/pkg/adapter/abs"
	"github.com/secmon-as-code/locust/pkg/domain/context/logging"
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
			EnvVars:     []string{"LOCUST_ENABLE_AZURE"},
			Destination: &x.enable,
		},

		&cli.StringFlag{
			Name:        "azure-tenant-id",
			Usage:       "Azure tenant ID",
			Category:    category,
			EnvVars:     []string{"LOCUST_AZURE_TENANT_ID"},
			Destination: &x.tenantID,
		},
		&cli.StringFlag{
			Name:        "azure-client-id",
			Usage:       "Azure client ID",
			Category:    category,
			EnvVars:     []string{"LOCUST_AZURE_CLIENT_ID"},
			Destination: &x.clientID,
		},
		&cli.StringFlag{
			Name:        "azure-client-secret",
			Usage:       "Azure client secret",
			Category:    category,
			EnvVars:     []string{"LOCUST_AZURE_CLIENT_SECRET"},
			Destination: &x.clientSecret,
		},
	}
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
