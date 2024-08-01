package main

import (
	"io"
	"log"
	"os"
)

func pack() {
	log.Printf("> cleanup...")
	if err := cleanup(); err != nil {
		log.Fatal(err)
	}

	install.Product.Name = makeInstall.Product.Name
	install.Target.Editable = makeInstall.Target.EditablePath
	install.Files.Type = makeInstall.Files.Type
	install.Files.Encrypt = makeInstall.Files.Encrypt
	if install.Files.Encrypt {
		initEncoderKey()
	}

	if n, err := makeInstall.ParseSplitSize(); err != nil {
		log.Fatal(err)
	} else {
		rawOutput.maxCount = n
	}

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
		out, err := output.Open(it)
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
