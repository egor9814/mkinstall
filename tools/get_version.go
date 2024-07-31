package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file := path.Join(wd, "version.ref")

	cmd := exec.Command("git", "describe", "--exact-match", "--tags")
	cmd.Dir = wd

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		err = os.WriteFile(file, []byte("0.0.0"), 0444)
	} else {
		err = os.WriteFile(file, out.Bytes(), 0444)
	}

	if err != nil {
		log.Fatal(err)
	}
}
