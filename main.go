package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path"
)

func exe() string {
	return path.Base(os.Args[0])
}

func help() {}

func version() {
	parseVersion()
	fmt.Printf("%s v%d.%d.%d%s\n", exe(), Version.Major, Version.Minor, Version.Patch, Version.Suffix)
}

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

		case "version":
			version()
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
