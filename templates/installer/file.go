package main

import (
	"bytes"
	"io"
	"os"
	"path"
)

type VirtualFile struct {
	Path  string
	Embed bool
}

type EmbedFiles map[string][]byte

var embedFiles EmbedFiles

type bytesReader struct {
	reader *bytes.Reader
}

func (r *bytesReader) Read(b []byte) (n int, err error) {
	if r.reader == nil {
		return 0, os.ErrClosed
	}
	return r.reader.Read(b)
}

func (r *bytesReader) Close() error {
	if r.reader == nil {
		return os.ErrClosed
	}
	r.reader = nil
	return nil
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
	if f.Embed {
		if embedFiles == nil {
			return nil, os.ErrNotExist
		}
		if data, ok := embedFiles[f.Path]; ok {
			return &bytesReader{
				reader: bytes.NewReader(data),
			}, nil
		} else {
			return nil, os.ErrNotExist
		}
	} else {
		return os.Open(f.Path)
	}
}

func (f *VirtualFile) IsValid() bool {
	return len(f.Path) != 0
}

func (f *VirtualFile) Size() (int, error) {
	if len(f.Path) == 0 {
		return 0, nil
	}
	if f.Embed {
		if embedFiles == nil {
			return 0, os.ErrNotExist
		}
		if data, ok := embedFiles[f.Path]; ok {
			return len(data), nil
		}
		return 0, os.ErrNotExist
	}
	info, err := os.Stat(f.Path)
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
