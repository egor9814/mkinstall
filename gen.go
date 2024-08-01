package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

//go:embed templates/installer/file.go
var template_file_go string

//go:embed templates/installer/go.mod_
var template_go_mod string

//go:embed templates/installer/go.sum_
var template_go_sum string

//go:embed templates/installer/input.go
var template_input_go string

//go:embed templates/installer/install.go
var template_install_go string

//go:embed templates/installer/lang.go
var template_lang_go string

//go:embed templates/installer/main.go
var template_main_go string

//go:embed templates/installer/raw_input.go
var template_raw_input_go string

//go:embed templates/installer/tar_input.go
var template_tar_input_go string

//go:embed version.go
var template_version_go string

func generate() error {
	write := func(name, data string) error {
		target := path.Join(workInstallerDir, name)
		if err := os.MkdirAll(path.Dir(target), 0700); err != nil {
			return err
		}
		return os.WriteFile(target, []byte(data), 0600)
	}

	log.Printf("> generating installer...\n")
	if err := write("file.go", template_file_go); err != nil {
		return err
	}

	if err := write("go.mod", template_go_mod); err != nil {
		return err
	}

	if err := write("go.mod", template_go_mod); err != nil {
		return err
	}

	if err := write("go.sum", template_go_sum); err != nil {
		return err
	}

	if err := write("input.go", template_input_go); err != nil {
		return err
	}

	if err := write("install.go", template_install_go); err != nil {
		return err
	}

	if err := write("lang.go", template_lang_go); err != nil {
		return err
	}

	if err := write("main.go", template_main_go); err != nil {
		return err
	}

	if err := write("raw_input.go", template_raw_input_go); err != nil {
		return err
	}

	if err := write("tar_input.go", template_tar_input_go); err != nil {
		return err
	}

	if err := write("version.go", template_version_go); err != nil {
		return err
	}

	if err := write("version.ref", version_ref); err != nil {
		return err
	}

	if install.Files.Embed {
		template := make([]string, 1, 2+2*len(install.Files.List))
		template[0] = `package main

import _ "embed"

func init() {
	embedFiles = make(EmbedFiles)`
		for i, it := range install.Files.List {
			template = append(template, "\tembedFiles[\""+it+"\"] = __embed_file_"+strconv.Itoa(i))
		}
		template = append(template, "}\n")

		for i, it := range install.Files.List {
			template = append(template, "//go:embed "+it+"\nvar __embed_file_"+strconv.Itoa(i)+" []byte\n")
		}

		if err := write("embed_files.go", strings.Join(template, "\n")); err != nil {
			return err
		}
	}

	pl := len(makeInstall.Target.Platforms)
	for i, it := range makeInstall.Target.Platforms {
		log.Printf("> [%d/%d] building installer for %s %s...\n", i+1, pl, it.Os, it.Arch)
		install.Target.Path = it.Path
		if data, err := install.Json(); err != nil {
			return err
		} else {
			if err := write("install.json", string(data)); err != nil {
				return err
			}
		}
		if err := buildInstaller(&it); err != nil {
			return err
		}
	}

	return nil
}

func buildInstaller(platform *TargetPlatform) error {
	target := path.Join(workOutputDir, "Setup_"+platform.Os+"_"+platform.Arch)
	if platform.Os == "windows" {
		target += ".exe"
	}
	cmd := exec.Command("go", "build", "-v", "-o", target)
	cmd.Dir = workInstallerDir
	cmd.Env = append(cmd.Env, "GOOS="+platform.Os, "GOARCH="+platform.Arch, "GOPATH="+goPath, "GOCACHE="+goCache)
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("> build error: `%v`; output:\n%s", err, out.String())
	}
	return err
}
