package main

import (
	_ "embed"
	"log"
	"os"
)

func help() {}

//go:embed mkinstall.json
var template_mkinstall_json string

func initProject() {
	if err := os.WriteFile("mkinstall.json", []byte(template_mkinstall_json), 0700); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initWorkDir()

	inputFile := "mkinstall.json"

	for i, l := 1, len(os.Args); i < l; i++ {
		it := os.Args[i]
		switch it {
		case "help":
			help()
			return

		case "init":
			initProject()
			return

		default:
			if _, err := os.Stat(it); err != nil {
				log.Fatalf("invaid argument or file %q", it)
			} else {
				inputFile = it
			}
		}
	}

	if _, err := os.Stat(inputFile); err != nil {
		log.Fatalf("file %q not found", inputFile)
	}

	if err := makeInstall.load(inputFile); err != nil {
		log.Fatal(err)
	}

	pack()
}
