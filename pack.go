package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func pack(noPack bool) {
	if noPack {
		install.Files = make([]string, 0, 64)
		err := filepath.Walk(workOutputDir, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasPrefix(info.Name(), "data-") && strings.HasSuffix(info.Name(), ".dat") {
				if f, err := filepath.Rel(workOutputDir, p); err != nil {
					return err
				} else {
					install.Files = append(install.Files, f)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		if len(install.Files) == 0 {
			noPack = false
			log.Printf("> warning: --no-pack flag provided, but dat files not found")
		}
	}
	if makeInstall.Files.Encrypt && noPack {
		keyFile := filepath.Join(workInstallerDir, "encoder.key")
		if data, err := os.ReadFile(keyFile); err != nil {
			log.Printf("> error: cannot read encryption key file %q\n", keyFile)
			log.Fatal(err)
		} else if len(data) != len(encoder.key) {
			log.Fatalf("> error: invalid encryption key file %q", keyFile)
		} else {
			copy(encoder.key[:], data)
			encoder.initialized = true
		}
	}

	log.Printf("> cleanup...")
	if err := cleanup(noPack); err != nil {
		log.Fatal(err)
	}

	install.ProductName = makeInstall.Product.Name
	install.TargetPathEditable = makeInstall.Target.EditablePath
	install.InputType = makeInstall.InputType()
	install.Decrypt = makeInstall.Files.Encrypt
	if install.Decrypt {
		initEncoderKey()
	}

	if n, err := makeInstall.ParseSplitSize(); err != nil {
		log.Fatal(err)
	} else {
		rawOutput.maxCount = n
	}

	install.Shortcuts = makeInstall.Shortcuts

	if !noPack {
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
	}

	if err := generate(); err != nil {
		log.Fatal(err)
	}

	log.Printf("> done!\n")
	log.Printf("> result in %q\n", workOutputDir)
}
