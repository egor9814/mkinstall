package main

import (
	"io"
	"math/rand"
)

type encoderType struct {
	wc    io.WriteCloser
	index int
	key   [128]byte
}

var encoder encoderType

func initEncoderKey() {
	for i := range encoder.key {
		n := rand.Uint64()
		for j := 1; j < 8; j++ {
			n ^= (n >> (8 * j) & 0xff)
		}
		encoder.key[i] = byte(n & 0xff)
	}
}

func (e *encoderType) Write(b []byte) (int, error) {
	if e.wc == nil {
		return 0, io.ErrClosedPipe
	}
	data := make([]byte, len(b))
	for i, it := range b {
		data[i] = it ^ e.key[e.index]
		e.index = (e.index + 1) % len(e.key)
	}
	return e.wc.Write(data)
}

func (e *encoderType) Close() (err error) {
	if e.wc != nil {
		err = e.wc.Close()
		e.wc = nil
	}
	return
}
