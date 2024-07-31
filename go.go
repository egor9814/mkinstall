package main

import (
	"log"
	"os"
)

var goPath, goCache string

func init() {
	found := false
	goPath, found = os.LookupEnv("GOPATH")
	if !found {
		log.Fatal("GOPATH not provided")
	}

	// var out bytes.Buffer
	// cmd := exec.Command("go", "env", "GOCACHE")
	// cmd.Stdout = &out
	// cmd.Env = append(cmd.Env, "GOPATH="+goPath)
	// if err := cmd.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	// goCache = strings.TrimSpace(out.String())
}
