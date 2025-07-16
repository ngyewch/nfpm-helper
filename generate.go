package main

import (
	"context"
	"github.com/ngyewch/nfpm-helper/generate"
	"github.com/ngyewch/nfpm-helper/utils"
	"github.com/urfave/cli/v3"
)

func doGenerate(ctx context.Context, cmd *cli.Command) error {
	packagers := cmd.StringSlice(packagersFlags.Name)
	outputDir := cmd.String(outputDirFlag.Name)

	var config generate.Config
	err := utils.LoadConfigurationFromFile("nfpm-helper.gen.yml", &config)
	if err != nil {
		return err
	}

	generator := &generate.Generator{
		Config:    config,
		Packagers: packagers,
		OutputDir: outputDir,
	}

	return generator.Generate(ctx)
}
