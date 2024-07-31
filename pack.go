package main

import (
	"io"
	"log"
	"os"
)

func pack() {
	install.Product.Name = makeInstall.Product.Name
	install.Target.Path = makeInstall.Target.Path
	install.Target.Editable = makeInstall.Target.Editable
	install.Files.Embed = false                 // TODO: support
	install.Files.Type = "zstd"                 // TODO: support
	rawOutput.maxCount = 8 * 1024 * 1024 * 1024 // TODO: support

	files, err := makeInstall.makeFiles()
	if err != nil {
		log.Fatal(err)
	}

	output, err := NewOutput()
	if err != nil {
		log.Fatal(err)
	}

	fl := len(files)
	for i, it := range files {
		log.Printf("> [%d/%d] packing %q...\n", i+1, fl, it)
		info, _ := os.Stat(it)
		out, err := output.Open(it, int(info.Size()))
		if err != nil {
			output.Close()
			log.Fatal(err)
		}

		in, err := os.Open(it)
		if err != nil {
			out.Close()
			output.Close()
			log.Fatal(err)
		}

		_, err = io.Copy(out, in)
		in.Close()
		out.Close()
		if err != nil {
			output.Close()
			log.Fatal(err)
		}
	}
	output.Close()

	err = generate()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("> done!\n")
	log.Printf("> result in %q\n", workOutputDir)
}
