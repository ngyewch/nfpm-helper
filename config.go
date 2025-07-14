package main

import (
	"github.com/goccy/go-yaml"
	"io"
	"os"
)

type Configuration struct {
	Name            string              `yaml:"name"`
	Download        DownloadBaseConfig  `yaml:"download"`
	StripComponents int                 `yaml:"strip_components"`
	Packaging       PackagingBaseConfig `yaml:"packaging"`
	Packagers       []string            `yaml:"packagers"`
	Outputs         []Output            `yaml:"outputs"`
}

type Output struct {
	Arch      string          `yaml:"arch"`
	Download  DownloadConfig  `yaml:"download"`
	Packaging PackagingConfig `yaml:"packaging"`
}

type DownloadBaseConfig struct {
	UrlTemplate string `yaml:"url_template"`
}

type DownloadConfig struct {
	UrlTemplate string            `yaml:"url_template"`
	Env         map[string]string `yaml:"env"`
}

type PackagingBaseConfig struct {
	FilenameTemplate string `yaml:"filename_template"`
}

type PackagingConfig struct {
	FilenameTemplate string            `yaml:"filename_template"`
	Env              map[string]string `yaml:"env"`
}

func LoadConfiguration(r io.Reader) (*Configuration, error) {
	var configuration Configuration
	decoder := yaml.NewDecoder(r, yaml.Strict())
	err := decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func LoadConfigurationFromFile(path string) (*Configuration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	return LoadConfiguration(f)
}
