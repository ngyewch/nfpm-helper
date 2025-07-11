package main

import (
	"github.com/goccy/go-yaml"
	"io"
	"os"
)

type Configuration struct {
	Name                   string   `yaml:"name"`
	DownloadUrlTemplate    string   `yaml:"download_url_template"`
	StripComponents        int      `yaml:"strip_components"`
	OutputFilenameTemplate string   `yaml:"output_filename_template"`
	Outputs                []Output `yaml:"outputs"`
	Packagers              []string `yaml:"packagers"`
}

type Output struct {
	Arch string `yaml:"arch"`
	Url  string `yaml:"url"`
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
