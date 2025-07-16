package main

import (
	"github.com/goccy/go-yaml"
	"io"
	"os"
)

func LoadConfiguration[T any](r io.Reader, v T) error {
	decoder := yaml.NewDecoder(r, yaml.Strict())
	return decoder.Decode(v)
}

func LoadConfigurationFromFile[T any](path string, v T) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	return LoadConfiguration(f, v)
}
