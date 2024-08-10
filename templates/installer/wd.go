package main

import (
	"os"
	"path/filepath"
)

var workDir string

func init() {
	workDir = filepath.Dir(os.Args[0])
}
