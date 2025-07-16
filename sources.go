package main

import (
	"github.com/boyter/gocodewalker"
	"os"
)

type Sources struct {
	Entries []Source `json:"entries"`
}

type Source struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"modTime"`
}

func getSources(path string) (*Sources, error) {
	var sources Sources

	fileListQueue := make(chan *gocodewalker.File, 100)
	fileWalker := gocodewalker.NewFileWalker(path, fileListQueue)
	go func() {
		_ = fileWalker.Start()
	}()

	for f := range fileListQueue {
		stat, err := os.Stat(f.Location)
		if err != nil {
			fileWalker.Terminate()
			return nil, err
		}
		entry := Source{
			Path:    f.Location,
			Size:    stat.Size(),
			ModTime: stat.ModTime().UnixMilli(),
		}
		sources.Entries = append(sources.Entries, entry)
	}

	return &sources, nil
}
