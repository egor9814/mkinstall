package main

import (
	"archive/tar"
	"io"
)

type tarInputImpl struct {
	reader *tar.Reader
	c      io.Closer
	i      int
}

func newTarInput(r io.Reader, c io.Closer) *tarInputImpl {
	return &tarInputImpl{
		reader: tar.NewReader(r),
		c:      c,
	}
}

func (i *tarInputImpl) Progress() ProgressStatus {
	return rawInput.Progress()
}

type tarReadCloser struct {
	i *tarInputImpl
}

func (r *tarReadCloser) Read(p []byte) (n int, err error) {
	n, err = r.i.reader.Read(p)
	if err == io.EOF && n != 0 {
		err = nil
	}
	return
}

func (*tarReadCloser) Close() error {
	return nil
}

func (i *tarInputImpl) Next() (InputFile, error) {
	header, err := i.reader.Next()
	if err == io.EOF {
		return InputFile{}, nil
	}
	if err != nil {
		return InputFile{}, err
	}
	i.i++
	return InputFile{
		Path: header.Name,
		Open: func() (io.ReadCloser, error) {
			return &tarReadCloser{
				i: i,
			}, nil
		},
	}, nil
}

func (i *tarInputImpl) Close() (err error) {
	i.reader = nil
	if i.c != nil {
		err = i.c.Close()
		i.c = nil
	}
	return
}
