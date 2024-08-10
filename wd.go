package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var workDir string
var workOutputDir string
var workInstallerDir string

var goPath, goCache, goTmp string

func init() {
	if wd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else {
		workDir = wd
	}
	workOutputDir = filepath.Join(workDir, "mkinstall-output")
	workInstallerDir = filepath.Join(workDir, "mkinstall-temp")

	goCache = filepath.Join(workInstallerDir, ".cache")
	goTmp = filepath.Join(workInstallerDir, ".tmp")

	found := false
	goPath, found = os.LookupEnv("GOPATH")
	if !found {
		cmd := exec.Command("go", "env", "GOPATH")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err == nil {
			goPath = strings.TrimSpace(out.String())
			info, err := os.Stat(goPath)
			if err == nil && info.IsDir() {
				found = true
			}
		}
	}
	if !found {
		log.Fatal("GOPATH not provided")
	}
}
