package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var workDir string
var workOutputDir string
var workInstallerDir string

var goPath, goCache string

func init() {
	if wd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else {
		workDir = wd
	}
	workOutputDir = path.Join(workDir, "mkinstall-output")
	workInstallerDir = path.Join(workDir, "mkinstall-temp")

	goCache = path.Join(workInstallerDir, ".cache")

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
