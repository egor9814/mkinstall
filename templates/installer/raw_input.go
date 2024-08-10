package main

import (
	"fmt"
	"io"
)

type rawInputImpl struct {
	index    int
	size     int64
	read     int64
	progress chan int64
}

var rawInput rawInputImpl

func (i *rawInputImpl) init() error {
	i.progress = make(chan int64)
	for it := i.NextFile(); it.IsValid(); it = i.NextFile() {
		size, err := it.Size()
		if err != nil {
			return err
		}
		i.size += size
	}
	i.index = 0
	return nil
}

func (i *rawInputImpl) All() int64 {
	return i.size
}

func (i *rawInputImpl) Current() int64 {
	return i.read
}

func (i *rawInputImpl) Chan() <-chan int64 {
	return i.progress
}

func (i *rawInputImpl) Progress() ProgressStatus {
	return i
}

func (i *rawInputImpl) NextFile() VirtualFile {
	if i.index >= len(install.Files) {
		return VirtualFile{}
	}
	i.index++
	return VirtualFile{
		Path: install.Files[i.index-1],
	}
}

func (i *rawInputImpl) Next() (InputFile, error) {
	file := i.NextFile()
	if !file.IsValid() {
		return InputFile{}, nil
	}
	return InputFile{
		Path: file.Path,
		Open: file.Open,
	}, nil
}

func (i *rawInputImpl) Close() error {
	i.index = len(install.Files)
	close(i.progress)
	i.progress = nil
	return nil
}

type rawReaderImpl struct {
	current io.ReadCloser
}

func (r *rawReaderImpl) next() error {
	if r.current != nil {
		if err := r.current.Close(); err != nil {
			return err
		}
	}
	if f := rawInput.NextFile(); f.IsValid() {
		rc, err := f.Open()
		if err == nil {
			r.current = rc
		}
		return err
	} else {
		r.current = nil
		return io.EOF
	}
}

func (r *rawReaderImpl) sendProgress() {
	defer func() {
		// supress send to closed channel
		_ = recover()
	}()
	if rawInput.progress != nil {
		rawInput.progress <- rawInput.read
	}
}

func (r *rawReaderImpl) Read(p []byte) (int, error) {
	if r.current == nil {
		if err := r.next(); err != nil {
			return 0, err
		}
	}
	n, err := r.current.Read(p)
	if err == io.EOF && n != 0 {
		err = nil
	}
	if err == nil {
		for n < len(p) {
			if err := r.next(); err != nil {
				return n, err
			}
			n2, err := r.current.Read(p[n:])
			n += n2
			if err == io.EOF {
				break
			}
			if err != nil {
				return n, err
			}
		}
	}
	rawInput.read += int64(n)
	r.sendProgress()
	return n, err
}

func (r *rawReaderImpl) Close() (err error) {
	if r.current != nil {
		err = r.current.Close()
		r.current = nil
	}
	if err2 := rawInput.Close(); err2 != nil {
		if err == nil {
			err = err2
		} else {
			err = fmt.Errorf("%v\n%v", err, err2)
		}
	}
	return
}

var rawReader rawReaderImpl
