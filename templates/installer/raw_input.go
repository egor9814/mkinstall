package main

import (
	"bytes"
	"io"
)

type rawInputImpl struct {
	allBytes  int
	index     int
	readCount int
}

func (i *rawInputImpl) init() error {
	if i.allBytes == 0 {
		for {
			if it := i.NextFile(); it.IsValid() {
				size, err := it.Size()
				if err != nil {
					return err
				}
				i.allBytes += size
			} else {
				break
			}
		}
	}
	i.index = 0
	return nil
}

var rawInput rawInputImpl

func (i *rawInputImpl) ProgressCurrent() int {
	return i.index
}

func (i *rawInputImpl) ProgressAll() int {
	return len(install.Files.List)
}

func (i *rawInputImpl) NextFile() VirtualFile {
	if i.index >= len(install.Files.List) {
		return VirtualFile{}
	}
	i.index++
	return VirtualFile{
		Path:  install.Files.List[i.index-1],
		Embed: install.Files.Embed,
	}
}

func (i *rawInputImpl) Next() (OutputFile, error) {
	file := i.NextFile()
	if !file.IsValid() {
		return OutputFile{}, nil
	}
	return OutputFile{
		Path: file.Path,
		Open: file.Open,
	}, nil
}

func (i *rawInputImpl) Close() error {
	i.index = len(install.Files.List)
	return nil
}

type rawReader struct {
	buf bytes.Buffer
}

func (r *rawReader) populate() error {
	if r.buf.Available() == 0 {
		if f := rawInput.NextFile(); f.IsValid() {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			io.Copy(&r.buf, rc)
		} else {
			return io.EOF
		}
	}
	return nil
}

func (r *rawReader) Read(p []byte) (int, error) {
	if err := r.populate(); err != nil {
		return 0, err
	}
	n, err := r.buf.Read(p)
	rawInput.readCount += n
	return n, err
}
