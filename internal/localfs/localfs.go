package localfs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type LocalFS interface {
	Open(path string) (fs.File, error)
	Create(path string) (*os.File, error)
}

type localFS struct {
	fdir string
	ffs  fs.FS
}

func NewLocalFS(fdir string) LocalFS {
	if !strings.HasSuffix(fdir, "/") {
		fdir += "/"
	}
	return &localFS{fdir: fdir, ffs: os.DirFS(fdir)}
}

func (l *localFS) Open(path string) (fs.File, error) {
	return l.ffs.Open(path)
}

func (l *localFS) Create(path string) (*os.File, error) {
	_, err := l.Open(path)
	if err == nil {
		return nil, os.ErrExist
	}

	f, err := os.Create(l.fdir + path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return f, nil
}
