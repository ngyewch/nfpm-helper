package main

import (
	"context"
	slogUtils "github.com/ngyewch/go-clibase/slog-utils"
	"github.com/urfave/cli/v3"
	"log/slog"
	"os"
)

var (
	version string

	log = slogUtils.GetLoggerForCurrentPackage()

	logLevelFlag = &cli.StringFlag{
		Name:     "log-level",
		Usage:    "log level",
		Category: "Logging",
		Value:    "INFO",
		Sources:  cli.EnvVars("LOG_LEVEL"),
		Action: func(ctx context.Context, cmd *cli.Command, s string) error {
			slogUtils.SetLevel(slogUtils.ToLevel(s))
			return nil
		},
	}

	outputDirFlag = &cli.StringFlag{
		Name:    "output-dir",
		Usage:   "output directory",
		Value:   "build",
		Sources: cli.EnvVars("OUTPUT_DIR"),
	}

	app = &cli.Command{
		Name:    "nfpm-helper",
		Usage:   "nfpm helper",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "build",
				Usage:  "build",
				Action: doBuild,
				Flags: []cli.Flag{
					outputDirFlag,
				},
			},
		},
		Flags: []cli.Flag{
			logLevelFlag,
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Error("error",
			slog.Any("err", err),
		)
	}
}
