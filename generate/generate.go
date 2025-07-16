package generate

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

		var absRepositoryDir string
		switch repository.Type {
		case "", "local":
			fmt.Printf("Using local repository %s\n", repository.Source)
			repositoryDir := os.ExpandEnv(repository.Source)
			absRepositoryDir, err = filepath.Abs(repositoryDir)
			if err != nil {
				return err
			}

		case "git":
			fmt.Printf("Cloning git repository %s\n", repository.Source)
			tempDir, err := os.MkdirTemp("", "nfpm-helper-")
			if err != nil {
				return err
			}
			defer func(path string) {
				_ = os.RemoveAll(path)
			}(tempDir)

			gitRepo, err := git.PlainClone(tempDir, false, &git.CloneOptions{
				URL: repository.Source,
			})
			if err != nil {
				return err
			}

			if repository.Version != "" {
				workTree, err := gitRepo.Worktree()
				if err != nil {
					return err
				}

				branch, err := gitRepo.Branch(repository.Version)
				if errors.Is(err, git.ErrBranchNotFound) {
					tag, err := gitRepo.Tag(repository.Version)
					if errors.Is(err, git.ErrTagNotFound) {
						fmt.Printf("Checking out commit %s\n", repository.Version)
						hash, err := gitRepo.ResolveRevision(plumbing.Revision(repository.Version))
						if err != nil {
							return err
						}
						err = workTree.Checkout(&git.CheckoutOptions{
							Hash:  *hash,
							Force: true,
						})
						if err != nil {
							return err
						}
					} else if err != nil {
						return err
					} else {
						fmt.Printf("Checking out tag %s\n", repository.Version)
						err = workTree.Checkout(&git.CheckoutOptions{
							Hash:  tag.Hash(),
							Force: true,
						})
						if err != nil {
							return err
						}
					}
				} else if err != nil {
					return err
				} else {
					fmt.Printf("Checking out branch %s\n", repository.Version)
					err = workTree.Checkout(&git.CheckoutOptions{
						Branch: branch.Merge,
						Force:  true,
					})
					if err != nil {
						return err
					}
				}
			}

			absRepositoryDir = tempDir

		default:
			return fmt.Errorf("repository type '%s' is not supported", repository.Type)
		}

		err = os.Chdir(absRepositoryDir)
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
