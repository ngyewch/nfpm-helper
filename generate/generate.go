package generate

import (
	"context"
	"fmt"
	"github.com/ngyewch/nfpm-helper/build"
	"github.com/ngyewch/nfpm-helper/utils"
	"os"
	"path/filepath"
)

type Generator struct {
	Config    Config
	Packagers []string
	OutputDir string
}

func (generator *Generator) Generate(ctx context.Context) error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, repository := range generator.Config.Repositories {
		err = os.Chdir(workingDirectory)
		if err != nil {
			return err
		}

		repositoryDir := os.ExpandEnv(repository.Dir)
		err = os.Chdir(repositoryDir)
		if err != nil {
			return err
		}

		absRepositoryDir, err := filepath.Abs(repositoryDir)
		if err != nil {
			return err
		}

		var index IndexConfig
		err = utils.LoadConfigurationFromFile("nfpm-helper.index.yml", &index)
		if err != nil {
			return err
		}

		for _, pkg := range repository.Packages {
			matched := false
			for _, pkgDef := range index.Packages {
				if pkgDef.Name == pkg.Name {
					fmt.Printf("Generating %s %s ...\n", pkg.Name, pkg.Version)

					pkgDir := os.ExpandEnv(pkgDef.Dir)

					err = os.Chdir(absRepositoryDir)
					if err != nil {
						return err
					}

					err = os.Chdir(pkgDir)
					if err != nil {
						return err
					}

					var buildConfig build.Config
					err = utils.LoadConfigurationFromFile("nfpm-helper.yml", &buildConfig)
					if err != nil {
						return err
					}

					packagers := generator.Packagers
					if len(pkg.Packagers) > 0 {
						packagers = pkg.Packagers
					}
					builder := &build.Builder{
						Config:    buildConfig,
						Version:   pkg.Version,
						Archs:     pkg.Archs,
						Packagers: packagers,
						OutputDir: generator.OutputDir,
					}
					err = builder.Build(ctx)
					if err != nil {
						return err
					}

					matched = true
					break
				}
			}
			if !matched {
				return fmt.Errorf("package '%s' not found", pkg.Name)
			}
		}
	}

	return nil
}
