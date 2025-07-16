package main

import (
	"context"
	"github.com/ngyewch/nfpm-helper/build"
	"github.com/ngyewch/nfpm-helper/utils"
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

	version := cmd.String(versionFlag.Name)
	archs := cmd.StringSlice(archsFlags.Name)
	packagers := cmd.StringSlice(packagersFlags.Name)
	outputDir := cmd.String(outputDirFlag.Name)

	var config build.Config
	err := utils.LoadConfigurationFromFile("nfpm-helper.yml", &config)
	if err != nil {
		return err
	}

	builder := &build.Builder{
		Config:    config,
		Version:   version,
		Archs:     archs,
		Packagers: packagers,
		OutputDir: outputDir,
	}

	return builder.Build(ctx)
}
