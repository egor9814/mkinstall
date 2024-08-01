package main

import (
	_ "embed"
	"io"
)

//go:embed encoder.key
var encoderKey []byte

type decoderType struct {
	rc    io.ReadCloser
	index int
}

var decoder decoderType

func (d *decoderType) Read(p []byte) (int, error) {
	if d.rc == nil {
		return 0, io.ErrClosedPipe
	}
	n, err := d.rc.Read(p)
	if err == nil {
		for i := range p[:n] {
			p[i] ^= encoderKey[d.index]
			d.index = (d.index + 1) % len(encoderKey)
		}
	}
	return n, err
}

func (d *decoderType) Close() (err error) {
	if d.rc != nil {
		err = d.rc.Close()
		d.rc = nil
	}
	return
}
