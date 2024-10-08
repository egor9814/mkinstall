package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed templates/installer/decoder.go
var template_decoder_go string

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

//go:embed templates/installer/shortcut_unix.go
var template_shortcut_unix_go string

//go:embed templates/installer/shortcut_windows.go
var template_shortcut_windows_go string

//go:embed templates/installer/shortcut.go
var template_shortcut_go string

//go:embed templates/installer/tar_input.go
var template_tar_input_go string

//go:embed templates/installer/wd.go
var template_wd_go string

//go:embed templates/installer/zstd.go
var template_zstd_go string

func generate() error {
	writeBytes := func(name string, data []byte) error {
		target := filepath.Join(workInstallerDir, name)
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0600)
	}
	write := func(name, data string) error {
		return writeBytes(name, []byte(data))
	}

	log.Printf("> generating installer...\n")
	if err := write("decoder.go", template_decoder_go); err != nil {
		return err
	}

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

	if err := write("shortcut_unix.go", template_shortcut_unix_go); err != nil {
		return err
	}

	if err := write("shortcut_windows.go", template_shortcut_windows_go); err != nil {
		return err
	}

	if err := write("shortcut.go", template_shortcut_go); err != nil {
		return err
	}

	if err := write("tar_input.go", template_tar_input_go); err != nil {
		return err
	}

	if err := write("wd.go", template_wd_go); err != nil {
		return err
	}

	if err := write("zstd.go", template_zstd_go); err != nil {
		return err
	}

	if err := write("version.go", VersionString()); err != nil {
		return err
	}

	if install.Decrypt {
		if err := writeBytes("encoder.key", encoder.key[:]); err != nil {
			return err
		}
	} else {
		if err := write("encoder.key", ""); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(goTmp, 0700); err != nil {
		return err
	}
	pl := len(makeInstall.Target.Platforms)
	for i, it := range makeInstall.Target.Platforms {
		log.Printf("> [%d/%d] building installer for %s %s...\n", i+1, pl, it.Os, it.Arch)
		install.TargetPath = it.Path
		if err := write("install_var.go", install.String()); err != nil {
			return err
		}
		if err := buildInstaller(&it); err != nil {
			return err
		}
	}

	return nil
}

func buildInstaller(platform *TargetPlatform) error {
	target := filepath.Join(workOutputDir, "Setup_"+platform.Os+"_"+platform.Arch)
	if platform.Os == "windows" {
		target += ".exe"
	}
	cmd := exec.Command("go", "build", "-v", "-o", target)
	cmd.Dir = workInstallerDir
	cmd.Env = append(cmd.Env, "GOOS="+platform.Os, "GOARCH="+platform.Arch, "GOPATH="+goPath, "GOCACHE="+goCache, "GOTMPDIR="+goTmp)
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("> build error: `%v`; output:\n%s", err, out.String())
	}
	return err
}
