package main

import (
	"os"
	"path"
)

var workDir string

func init() {
	workDir = path.Dir(os.Args[0])
}
