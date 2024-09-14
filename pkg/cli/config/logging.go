package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

type Logging struct {
	level  string
	format string
	output string
}

func (x *Logging) Flags() []cli.Flag {
	const category = "Logging"

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Log level",
			Category:    category,
			EnvVars:     []string{"NYDUS_LOG_LEVEL"},
			Destination: &x.level,
			Value:       "info",
		},
		&cli.StringFlag{
			Name:        "log-format",
			Usage:       "Log format",
			Category:    category,
			EnvVars:     []string{"NYDUS_LOG_FORMAT"},
			Destination: &x.format,
			Value:       "json",
		},
		&cli.StringFlag{
			Name:        "log-output",
			Usage:       "Log output, stdout, stderr, or file path. Default is stdout",
			Category:    category,
			EnvVars:     []string{"NYDUS_LOG_OUTPUT"},
			Destination: &x.output,
			Value:       "-",
		},
	}
}

func (x *Logging) NewLogger() (*slog.Logger, error) {
	// Log output
	var output io.Writer
	switch x.output {
	case "stdout", "-":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		f, err := os.OpenFile(filepath.Clean(x.output), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil, goerr.Wrap(err, "failed to open log file").With("path", x.output)
		}
		output = f
	}

	// Log level
	levelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	level, ok := levelMap[x.level]
	if !ok {
		return nil, goerr.New("invalid log level").With("level", x.level)
	}

	// Log format
	var handler slog.Handler
	switch x.format {
	case "console":
		handler = clog.New(
			clog.WithWriter(output),
			clog.WithLevel(level),
			clog.WithSource(true),

			// clog.WithTimeFmt("2006-01-02 15:04:05"),
			clog.WithColorMap(&clog.ColorMap{
				Level: map[slog.Level]*color.Color{
					slog.LevelDebug: color.New(color.FgGreen, color.Bold),
					slog.LevelInfo:  color.New(color.FgCyan, color.Bold),
					slog.LevelWarn:  color.New(color.FgYellow, color.Bold),
					slog.LevelError: color.New(color.FgRed, color.Bold),
				},
				LevelDefault: color.New(color.FgBlue, color.Bold),
				Time:         color.New(color.FgWhite),
				Message:      color.New(color.FgHiWhite),
				AttrKey:      color.New(color.FgHiCyan),
				AttrValue:    color.New(color.FgHiWhite),
			}),
		)
	case "json":
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})

	default:
		return nil, goerr.New("invalid log format").With("format", x.format)
	}

	return slog.New(handler), nil
}
