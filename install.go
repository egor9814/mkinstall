package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"unicode"
)

type InstallInfo struct {
	Product struct {
		Name string `json:"name"`
	} `json:"product"`
	Target struct {
		Path     string `json:"path"`
		Editable bool   `json:"editable"`
	} `json:"target"`
	Files struct {
		Embed bool     `json:"embed"`
		Type  string   `json:"type"`
		List  []string `json:"list"`
	} `json:"files"`
}

var install InstallInfo

func (ii *InstallInfo) Json() ([]byte, error) {
	return json.Marshal(ii)
}

type TargetPlatform struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
	Path string `json:"path"`
}

type MakeInstallInfo struct {
	Product struct {
		Name string `json:"name"`
	} `json:"product"`
	Target struct {
		EditablePath bool             `json:"editable_path"`
		Platforms    []TargetPlatform `json:"platforms"`
	} `json:"target"`
	Files struct {
		Embed   bool     `json:"embed"`
		Type    string   `json:"type"`
		Split   string   `json:"split"`
		Include []string `json:"include"`
		Exclude []string `json:"exclude"`
	} `json:"files"`
}

var makeInstall MakeInstallInfo

func (ii *MakeInstallInfo) load(name string) error {
	if data, err := os.ReadFile(name); err != nil {
		return err
	} else {
		return json.Unmarshal(data, ii)
	}
}

func (ii *MakeInstallInfo) makeFiles() (list []string, err error) {
	if len(ii.Files.Include) == 0 {
		ii.Files.Include = append(ii.Files.Include, ".")
	}
	resolve := func(list []string, p string) ([]string, error) {
		if !path.IsAbs(p) {
			p = path.Join(workDir, p)
		}
		if _, err := os.Stat(p); err != nil {
			return list, os.ErrNotExist
		}
		err := filepath.Walk(p, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				list = append(list, p)
			}
			return nil
		})
		return list, err
	}
	included := make([]string, 0, 64)
	for _, it := range ii.Files.Include {
		included, err = resolve(included, it)
		if err != nil {
			return
		}
	}
	excluded := make([]string, 0, 64)
	ii.Files.Exclude = append(ii.Files.Exclude, "mkinstall.json", ".mkinstall", "mkinstall-output")
	for _, it := range ii.Files.Exclude {
		excluded, _ = resolve(excluded, it)
	}
	list = make([]string, 0, len(included))
loop:
	for _, it := range included {
		for _, ex := range excluded {
			if it == ex {
				continue loop
			}
		}
		list = append(list, it)
	}
	return
}

func (ii *MakeInstallInfo) Json() ([]byte, error) {
	return json.MarshalIndent(ii, "", "  ")
}

func (ii *MakeInstallInfo) ParseSplitSize() (int, error) {
	maximum := uint64(9223372036854775807) // int64.max
	if ii.Files.Embed {
		return int(maximum), nil
	}
	i := 0
	s := []rune(ii.Files.Split)
	l := len(s)
	for ; i < l; i++ {
		if !unicode.IsDigit(s[i]) {
			break
		}
	}
	if i == 0 {
		return 0, errors.New("files.split has invalid format")
	}
	n, err := strconv.ParseUint(string(s[:i]), 10, 64)
	if err == nil && i != l {
		for _, it := range "KMGT" {
			n *= 1024
			if it == s[i] {
				break
			}
		}
	}
	if err == nil && n == 0 {
		n = maximum
	}
	return int(n), err
}
