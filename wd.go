package main

import (
	"log"
	"os"
	"path"
)

var workDir string
var workOutputDir string
var workInstallerDir string

func initWorkDir() {
	if wd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else {
		workDir = wd
	}
	workOutputDir = path.Join(workDir, "mkinstall-output")
	workInstallerDir = path.Join(workDir, ".mkinstall")

	goCache = path.Join(workInstallerDir, ".cache")
}
