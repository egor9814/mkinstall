package main

import (
	"io"
	"os"
	"path"
	"strconv"
)

type rawOutputImpl struct {
	maxCount int
}

func (o *rawOutputImpl) Open(name string, size int) (io.WriteCloser, error) {
	target := path.Join(workOutputDir, name)
	if err := os.MkdirAll(path.Dir(target), 0700); err != nil {
		return nil, err
	}
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return nil, err
	}
	if install.Files.List == nil {
		install.Files.List = make([]string, 0, 64)
	}
	install.Files.List = append(install.Files.List, name)
	return out, nil
}

func (o *rawOutputImpl) Close() error {
	return nil
}

var rawOutput rawOutputImpl

type rawOutputWriter struct {
	current io.WriteCloser
	index   int
	count   int
}

func (w *rawOutputWriter) open() (err error) {
	if w.current != nil {
		err = w.current.Close()
	}
	if err == nil {
		w.index++
		w.current, err = rawOutput.Open("data-"+strconv.Itoa(w.index)+".dat", rawOutput.maxCount)
		if err == nil {
			w.count = 0
		} else {
			w.index--
		}
	}
	return
}

func (w *rawOutputWriter) Write(b []byte) (n int, err error) {
	if w.current == nil {
		err = w.open()
		if err != nil {
			return
		}
	}
	rem := rawOutput.maxCount - w.count
	for {
		if l := len(b); l > rem {
			n, err = w.current.Write(b[:rem])
			if err == nil && n < rem {
				err = io.ErrShortWrite
			}
			if err == nil {
				err = w.open()
			}
			if err != nil {
				return
			}
			b = b[rem:]
			rem = rawOutput.maxCount
		} else {
			rem = l
			break
		}
	}
	n, err = w.current.Write(b)
	if err == nil && n != rem {
		err = io.ErrShortWrite
	}
	if err == nil {
		w.count += n
	}
	return
}

func (w *rawOutputWriter) Close() (err error) {
	if w.current != nil {
		err = w.current.Close()
		w.current = nil
		w.index = 0
	}
	return
}
