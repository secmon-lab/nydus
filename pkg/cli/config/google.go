package config

import (
	"log/slog"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/nydus/pkg/adapter/gcs"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/option"
)

type GoogleCloudStorage struct {
	enable         bool
	credentialFile string
}

func (x *GoogleCloudStorage) Flags() []cli.Flag {
	const category = "Google Cloud Storage"

	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "enable-gcs",
			Usage:       "Enable Google Cloud Storage",
			Category:    category,
			EnvVars:     []string{"NYDUS_ENABLE_GCS"},
			Destination: &x.enable,
		},
		&cli.StringFlag{
			Name:        "gcs-credential-file",
			Usage:       "Google Cloud Storage credential file",
			Category:    category,
			EnvVars:     []string{"NYDUS_GCS_CREDENTIAL_FILE"},
			Destination: &x.credentialFile,
		},
	}
}

func (x GoogleCloudStorage) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Bool("enable", x.enable),
		slog.String("credentialFile", x.credentialFile),
	)
}

func (x *GoogleCloudStorage) NewClient() (*gcs.Client, error) {
	if !x.enable {
		return nil, nil
	}

	var options []gcs.Option
	if x.credentialFile != "" {
		options = append(options, gcs.WithGoogleAPIOption(option.WithCredentialsFile(x.credentialFile)))
	}

	client, err := gcs.New(options...)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create Google Cloud Storage client")
	}

	return client, nil
}
