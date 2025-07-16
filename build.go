package main

import (
	"context"
	"github.com/ngyewch/nfpm-helper/build"
	"github.com/urfave/cli/v3"
	"os"
)

func doBuild(ctx context.Context, cmd *cli.Command) error {
	workingDir := cmd.Args().First()
	if workingDir != "" {
		err := os.Chdir(workingDir)
		if err != nil {
			return err
		}
	}

	version := cmd.String(versionFlag.String())
	packagers := cmd.StringSlice(packagersFlags.Name)
	outputDir := cmd.String(outputDirFlag.Name)

	var config build.Configuration
	err := LoadConfigurationFromFile("nfpm-helper.yml", &config)
	if err != nil {
		return err
	}

	builder := &build.Builder{
		Config:    config,
		Version:   version,
		Packagers: packagers,
		OutputDir: outputDir,
	}

	return builder.Build(ctx)
}
