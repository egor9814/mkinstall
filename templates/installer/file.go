package main

import (
	"io"
	"os"
	"path/filepath"
)

type VirtualFile struct {
	Path string
}

func (f *VirtualFile) Open() (rc io.ReadCloser, err error) {
	if len(f.Path) == 0 {
		return nil, nil
	}
	rc, err = os.Open(filepath.Join(workDir, f.Path))
	if err == nil && install.Decrypt {
		decoder.rc = rc
		rc = &decoder
	}
	return
}

func (f *VirtualFile) IsValid() bool {
	return len(f.Path) != 0
}

func (f *VirtualFile) Create() (*os.File, error) {
	p := filepath.ToSlash(f.Path)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0600)
}

func (f *VirtualFile) Size() (int64, error) {
	if len(f.Path) == 0 {
		return 0, os.ErrNotExist
	}
	if info, err := os.Stat(filepath.Join(workDir, f.Path)); err != nil {
		return 0, err
	} else {
		return info.Size(), err
	}
}
