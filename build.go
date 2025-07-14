package main

import (
	"context"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/codeclysm/extract/v4"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
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

func doBuild(ctx context.Context, cmd *cli.Command) error {
	version := os.Getenv("VERSION")
	if version == "" {
		return fmt.Errorf("VERSION environment variable is not set")
	}

	config, err := LoadConfigurationFromFile("nfpm-helper.yml")
	if err != nil {
		return err
	}

	if len(config.Packagers) == 0 {
		return fmt.Errorf("no packagers specified")
	}

	for _, output := range config.Outputs {
		err = doOutput(ctx, config, output)
		if err != nil {
			return err
		}
	}

	return nil
}

func doOutput(ctx context.Context, config *Configuration, output Output) error {
	fmt.Printf("Packaging %s\n", output.Arch)

	customExpander := NewCustomExpander()
	customExpander.SetVar("NAME", config.Name)
	customExpander.SetVar("VERSION", os.Getenv("VERSION"))
	customExpander.SetVar("ARCH", output.Arch)

	downloadUrlTemplate := output.Download.UrlTemplate
	if downloadUrlTemplate == "" {
		downloadUrlTemplate = config.Download.UrlTemplate
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

	if config.StripComponents > 0 {
		renamer = func(p string) string {
			parts := strings.Split(p, "/")
			return strings.Join(parts[config.StripComponents:], "/")
		}
	}

	err = extract.Archive(ctx, f, tempDir, renamer)
	if err != nil {
		return err
	}
	customExpander.SetVar("ARCHIVE_DIR", tempDir)

	outputFilenameTemplate := output.Packaging.FilenameTemplate
	if outputFilenameTemplate == "" {
		outputFilenameTemplate = config.Packaging.FilenameTemplate
	}
	if outputFilenameTemplate == "" {
		outputFilenameTemplate = defaultOutputFilenameTemplate
	}

	packagingCustomExpander := customExpander.Clone()
	packagingCustomExpander.SetVars(output.Packaging.Env)
	outputFilename := packagingCustomExpander.Expand(outputFilenameTemplate)

	for _, packager := range config.Packagers {
		outputPath := filepath.Join("build", outputFilename+"."+packager)
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
	cachePath := filepath.Join(xdg.CacheHome, "nfpm-helper", u.Host, u.Path)
	_, err = os.Stat(cachePath)
	if os.IsNotExist(err) {
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

		fmt.Printf("Downloading %s\n", releaseUrl)
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
