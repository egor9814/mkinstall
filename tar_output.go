package main

import (
	"archive/tar"
	"io"
	"path/filepath"
)

type tarOutputImpl struct {
	w0 io.WriteCloser
	w  *tar.Writer
	f  func() error
}

func newTarOutput(w io.WriteCloser, flush func() error) *tarOutputImpl {
	return &tarOutputImpl{
		w0: w,
		w:  tar.NewWriter(w),
		f:  flush,
	}
}

type tarWriter struct {
	i *tarOutputImpl
}

func (w *tarWriter) Write(b []byte) (int, error) {
	return w.i.w.Write(b)
}

func (w *tarWriter) Close() error {
	// if w.i.f != nil {
	// 	return w.i.f()
	// }
	return nil
}

func (o *tarOutputImpl) Open(name string, size int) (io.WriteCloser, error) {
	realName, err := filepath.Rel(workDir, name)
	if err != nil {
		return nil, err
	}
	header := tar.Header{
		Name: realName,
		Mode: 0600,
		Size: int64(size),
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
