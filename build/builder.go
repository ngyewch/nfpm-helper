package build

import (
	"context"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/codeclysm/extract/v4"
	"github.com/ngyewch/nfpm-helper/utils"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	defaultOutputFilenameTemplate = "${NAME}_${VERSION}_${ARCH}"
)

type Builder struct {
	Config    Config
	Version   string
	Archs     []string
	Packagers []string
	OutputDir string
}

func (builder *Builder) Build(ctx context.Context) error {
	if len(builder.Archs) == 0 {
		return fmt.Errorf("no archs specified")
	}
	for _, arch := range builder.Archs {
		matched := false
		for _, output := range builder.Config.Outputs {
			if output.Arch == arch {
				err := builder.buildOutput(ctx, output)
				if err != nil {
					return err
				}
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("arch %s is not supported", arch)
		}
	}
	return nil
}

func (builder *Builder) buildOutput(ctx context.Context, output Output) error {
	fmt.Printf("Packaging %s\n", output.Arch)

	customExpander := utils.NewCustomExpander()
	customExpander.SetVar("NAME", builder.Config.Name)
	customExpander.SetVar("VERSION", builder.Version)
	customExpander.SetVar("ARCH", output.Arch)

	downloadUrlTemplate := output.Download.UrlTemplate
	if downloadUrlTemplate == "" {
		downloadUrlTemplate = builder.Config.Download.UrlTemplate
	}
	downloadCustomExpander := customExpander.Clone()
	downloadCustomExpander.SetVars(output.Download.Env)
	downloadUrl := downloadCustomExpander.Expand(downloadUrlTemplate)
	cachePath, err := download(ctx, downloadUrl)
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "nfpm-helper-*")
	if err != nil {
		return err
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tempDir)

	f, err := os.Open(cachePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var renamer extract.Renamer

	if builder.Config.StripComponents > 0 {
		renamer = func(p string) string {
			parts := strings.Split(p, "/")
			return strings.Join(parts[builder.Config.StripComponents:], "/")
		}
	}

	err = extract.Archive(ctx, f, tempDir, renamer)
	if err != nil {
		return err
	}
	customExpander.SetVar("ARCHIVE_DIR", tempDir)

	outputFilenameTemplate := output.Packaging.FilenameTemplate
	if outputFilenameTemplate == "" {
		outputFilenameTemplate = builder.Config.Packaging.FilenameTemplate
	}
	if outputFilenameTemplate == "" {
		outputFilenameTemplate = defaultOutputFilenameTemplate
	}

	packagingCustomExpander := customExpander.Clone()
	packagingCustomExpander.SetVars(output.Packaging.Env)
	outputFilename := packagingCustomExpander.Expand(outputFilenameTemplate)

	for _, packager := range builder.Packagers {
		outputPath := filepath.Join(builder.OutputDir, outputFilename+"."+packager)
		err = os.MkdirAll(filepath.Dir(outputPath), 0755)
		if err != nil {
			return err
		}

		var envList []string
		envList = append(envList, os.Environ()...)
		envList = append(envList, packagingCustomExpander.Environ()...)

		cmd := exec.CommandContext(ctx, "nfpm", "package",
			"--packager", packager,
			"--target", outputPath)
		cmd.Env = envList
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	fmt.Println()
	return nil
}

func download(ctx context.Context, releaseUrl string) (string, error) {
	u, err := url.Parse(releaseUrl)
	if err != nil {
		return "", err
	}
	cachePath := filepath.Join(xdg.CacheHome, "nfpm-helper", "downloads", u.Host, u.Path)
	_, err = os.Stat(cachePath)
	if os.IsNotExist(err) {
		fmt.Printf("Downloading %s\n", releaseUrl)
		httpResponse, err := http.Get(releaseUrl)
		if err != nil {
			return "", err
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(httpResponse.Body)

		if httpResponse.StatusCode != 200 {
			return "", fmt.Errorf(httpResponse.Status)
		}

		err = os.MkdirAll(filepath.Dir(cachePath), 0755)
		if err != nil {
			return "", err
		}

		f, err := os.OpenFile(cachePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return "", err
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		bar := progressbar.DefaultBytes(
			httpResponse.ContentLength,
			"",
		)
		_, err = io.Copy(io.MultiWriter(f, bar), httpResponse.Body)
		if err != nil {
			_ = os.Remove(cachePath)
			return "", err
		}

		return cachePath, nil
	} else if err != nil {
		return "", err
	} else {
		return cachePath, nil
	}
}
