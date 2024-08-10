package main

import (
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

type IOutput interface {
	io.Closer
	Open(name string) (io.WriteCloser, error)
}

func NewOutput() (IOutput, error) {
	switch makeInstall.Files.Type {
	case "raw":
		return &rawOutput, nil

	case "tar":
		return newTarOutput(&rawOutputWriter{}), nil

	case "zstd":
		// TODO: support zstd compression level
		coder, err := zstd.NewWriter(&rawOutputWriter{}, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
		if err != nil {
			return nil, err
		}
		return newTarOutput(coder), nil

	default:
		return nil, fmt.Errorf("unsupported files.type %q (supported: raw, tar, zstd)", makeInstall.Files.Type)
	}
}
