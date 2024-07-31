package main

import (
	"archive/tar"
	"io"
)

type tarInputImpl struct {
	reader *tar.Reader
}

func newTarInput(r io.Reader) *tarInputImpl {
	return &tarInputImpl{
		reader: tar.NewReader(r),
	}
}

func (i *tarInputImpl) ProgressCurrent() int {
	return rawInput.readCount
}

func (i *tarInputImpl) ProgressAll() int {
	return rawInput.allBytes
}

type tarReadCloser struct {
	i *tarInputImpl
	r io.Reader
}

func (r *tarReadCloser) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err == io.EOF && n != 0 {
		err = nil
	}
	return
}

func (*tarReadCloser) Close() error {
	return nil
}

func (i *tarInputImpl) Next() (OutputFile, error) {
	header, err := i.reader.Next()
	if err == io.EOF {
		return OutputFile{}, nil
	}
	if err != nil {
		return OutputFile{}, err
	}
	return OutputFile{
		Path: header.Name,
		Open: func() (io.ReadCloser, error) {
			return &tarReadCloser{
				i: i,
				r: i.reader,
			}, nil
		},
	}, nil
}

func (i *tarInputImpl) Close() error {
	i.reader = nil
	return nil
}
