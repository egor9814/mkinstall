package main

import (
	"io"
	"os"
	"path"
)

type VirtualFile struct {
	Path string
}

func (f *VirtualFile) Open() (rc io.ReadCloser, err error) {
	defer func() {
		if rc != nil && install.Files.Encrypt {
			decoder.rc = rc
			rc = &decoder
		}
	}()
	if len(f.Path) == 0 {
		return nil, nil
	}
	return os.Open(path.Join(workDir, f.Path))
}

func (f *VirtualFile) IsValid() bool {
	return len(f.Path) != 0
}

func (f *VirtualFile) Size() (int, error) {
	if len(f.Path) == 0 {
		return 0, nil
	}
	info, err := os.Stat(path.Join(workDir, f.Path))
	if err != nil {
		return 0, err
	}
	return int(info.Size()), nil
}

func (f *VirtualFile) Create() (*os.File, error) {
	dir := path.Dir(f.Path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return os.OpenFile(f.Path, os.O_CREATE|os.O_WRONLY, 0600)
}
