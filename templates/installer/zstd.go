package main

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

var zstdInputInfo struct {
	maxMem uint64
}

func newZstdInput(r io.Reader, c io.Closer) (IInput, error) {
	input, err := zstd.NewReader(
		r,
		zstd.WithDecoderConcurrency(0),
		zstd.WithDecoderLowmem(false),
		zstd.WithDecoderMaxMemory(zstdInputInfo.maxMem),
	)
	if err != nil {
		return nil, err
	}
	return newTarInput(input, c), nil
}
