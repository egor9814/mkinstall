package main

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

type OutputFile struct {
	Path string
	Open func() (io.ReadCloser, error)
}

func (f *OutputFile) IsValid() bool {
	return len(f.Path) != 0
}

type IInput interface {
	io.Closer
	ProgressCurrent() int
	ProgressAll() int
	Next() (OutputFile, error)
}

func NewInput() (IInput, error) {
	if err := rawInput.init(); err != nil {
		return nil, err
	}
	switch install.Files.Type {
	case "raw":
		return &rawInput, nil

	case "tar":
		return newTarInput(&rawReader{}), nil

	case "zstd":
		input, err := zstd.NewReader(&rawReader{})
		if err != nil {
			return nil, err
		}
		return newTarInput(input), nil

	default:
		panic("unreachable")
	}
}
