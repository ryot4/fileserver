package main

import (
	"net/http"
	"os"
	"strings"
)

type dotFileHidingFs struct {
	http.FileSystem
}

func (fs dotFileHidingFs) Open(name string) (http.File, error) {
	for _, part := range strings.Split(name, "/") {
		if strings.HasPrefix(part, ".") {
			return nil, os.ErrNotExist
		}
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return dotFileHidingFile{file}, nil
}

type dotFileHidingFile struct {
	http.File
}

func (dir dotFileHidingFile) Readdir(count int) (filtered []os.FileInfo, err error) {
	files, err := dir.File.Readdir(count)
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			filtered = append(filtered, f)
		}
	}
	return
}
