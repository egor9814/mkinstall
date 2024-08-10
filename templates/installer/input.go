package main

import (
	"io"
)

type InputType uint

const (
	RawInput InputType = iota
	TarInput
	ZstdInput
)

type InputFile struct {
	Path string
	Open func() (io.ReadCloser, error)
}

func (f *InputFile) IsValid() bool {
	return len(f.Path) != 0
}

type ProgressStatus interface {
	All() int64
	Current() int64
	Chan() <-chan int64
}

type IInput interface {
	io.Closer
	Progress() ProgressStatus
	Next() (InputFile, error)
}

func (t InputType) Open() (IInput, error) {
	if err := rawInput.init(); err != nil {
		return nil, err
	}
	switch t {
	case RawInput:
		return &rawInput, nil

	case TarInput:
		return newTarInput(&rawReader, &rawReader), nil

	case ZstdInput:
		return newZstdInput(&rawReader, &rawReader)

	default:
		panic("unreachable")
	}
}
