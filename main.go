package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func exe() string {
	return path.Base(os.Args[0])
}

func help() {
	fmt.Println("Make Install - Utilty for making installers")
	fmt.Printf("Usage: %s [COMMAND|FILE]\n", exe())
	fmt.Println("File by default is 'mkinstall.json'")
	fmt.Println("Commands:")
	fmt.Println(" help           - Print this help")
	fmt.Println(" version        - Print utility version")
	fmt.Println(" init           - Create mkinstall.json file in current directory")
}

func version() {
	parseVersion()
	fmt.Printf("%s v%d.%d.%d%s\n", exe(), Version.Major, Version.Minor, Version.Patch, Version.Suffix)
}

func initProject() {
	makeInstall.Product.Name = path.Base(workDir)
	makeInstall.Target.EditablePath = true
	platformOs := []string{
		"windows", "linux", "darwin", "android", "freebsd", "netbsd", "openbsd", "dragonfly", "plan9", "nacl",
	}
	platformArch := []string{
		"amd64", "arm", "amd64p32", "386",
	}
	makeInstall.Target.Platforms = make([]TargetPlatform, len(platformOs)*len(platformArch))
	for i, os := range platformOs {
		for j, arch := range platformArch {
			makeInstall.Target.Platforms[i*len(platformArch)+j] = TargetPlatform{
				Os:   os,
				Arch: arch,
				Path: "$HOME/",
			}
		}
	}
	makeInstall.Files.Type = "zstd"
	makeInstall.Files.Split = "8G"
	makeInstall.Files.Encrypt = false
	makeInstall.Files.Include = make([]string, 0)
	makeInstall.Files.Exclude = make([]string, 0)
	data, err := makeInstall.Json()
	if err == nil {
		err = os.WriteFile("mkinstall.json", data, 0600)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
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
