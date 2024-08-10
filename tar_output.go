package main

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

type tarOutputImpl struct {
	w0 io.WriteCloser
	w  *tar.Writer
}

func newTarOutput(w io.WriteCloser) *tarOutputImpl {
	return &tarOutputImpl{
		w0: w,
		w:  tar.NewWriter(w),
	}
}

type tarWriter struct {
	i *tarOutputImpl
}

func (w *tarWriter) Write(b []byte) (int, error) {
	return w.i.w.Write(b)
}

func (w *tarWriter) Close() error {
	return nil
}

func (o *tarOutputImpl) Open(name string) (io.WriteCloser, error) {
	realName, err := filepath.Rel(workDir, name)
	if err != nil {
		return nil, err
	}
	info, _ := os.Stat(name)
	header := tar.Header{
		Name: filepath.ToSlash(realName),
		Mode: 0600,
		Size: info.Size(),
	}
	if err := o.w.WriteHeader(&header); err != nil {
		return nil, err
	}
	return &tarWriter{
		i: o,
	}, nil
}

func (o *tarOutputImpl) Close() error {
	if err := o.w.Close(); err != nil {
		return err
	}
	return o.w0.Close()
}
