package main

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

type IOutput interface {
	io.Closer
	Open(name string, size int) (io.WriteCloser, error)
}

func NewOutput() (IOutput, error) {
	writer := &rawOutputWriter{}
	coder, err := zstd.NewWriter(writer, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return nil, err
	}
	return newTarOutput(coder), nil
}
