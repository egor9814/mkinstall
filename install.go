package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type InputType uint

const (
	RawInput InputType = iota
	TarInput
	ZstdInput
)

type InstallInfo struct {
	ProductName        string
	TargetPath         string
	TargetPathEditable bool
	InputType          string
	Decrypt            bool
	Files              []string
	Shortcuts          Shortcuts
}

var install InstallInfo

func (ii *InstallInfo) String() string {
	return fmt.Sprintf(`package main

var install = InstallInfo{
	ProductName:        %q,
	TargetPath:         %q,
	TargetPathEditable: %v,
	InputType:          %s,
	Decrypt:            %v,
	Files:              []string{"%s"},
	Shortcuts:          %s,
}
`, ii.ProductName, ii.TargetPath, ii.TargetPathEditable, ii.InputType, ii.Decrypt, strings.Join(ii.Files, `","`), ii.Shortcuts.String())
}

type TargetPlatform struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
	Path string `json:"path"`
}

type ShortcutUnix struct {
	Icon       *string  `json:"icon"`
	Categories []string `json:"categories"`
}

type Shortcut struct {
	Name         *string      `json:"name"`
	Target       string       `json:"target"`
	Arguments    []string     `json:"args"`
	ShortcutUnix ShortcutUnix `json:"unix"`
}

func (s *Shortcut) String() string {
	var name string
	if s.Name != nil {
		name = *s.Name
	} else {
		name = install.ProductName
	}
	var icon string
	if s.ShortcutUnix.Icon != nil {
		icon = *s.ShortcutUnix.Icon
	}
	return fmt.Sprintf(
		`{Name: %q, Target: %q, Arguments: []string{"%s"}, Icon: %q, Categories: []string{"%s"}}`,
		name,
		s.Target,
		strings.Join(s.Arguments, `", "`),
		icon,
		strings.Join(s.ShortcutUnix.Categories, `", "`),
	)
}

type Shortcuts []Shortcut

func (s *Shortcuts) String() string {
	l := make([]string, len(*s))
	for i, it := range *s {
		l[i] = it.String()
	}
	return fmt.Sprintf(
		`[]Shortcut{%s}`,
		strings.Join(l, ", "),
	)
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
		Type    string   `json:"type"`
		Split   string   `json:"split"`
		Encrypt bool     `json:"encrypt"`
		Include []string `json:"include"`
		Exclude []string `json:"exclude"`
	} `json:"files"`
	Shortcuts Shortcuts `json:"shortcuts"`
}

var makeInstall MakeInstallInfo

func (ii *MakeInstallInfo) load(name string) error {
	if data, err := os.ReadFile(name); err != nil {
		return err
	} else {
		err = json.Unmarshal(data, ii)
		if err == nil && len(makeInstall.Target.Platforms) == 0 {
			return errors.New("no target platforms")
		}
		return err
	}
}

func (ii *MakeInstallInfo) makeFiles() (list []string, err error) {
	if len(ii.Files.Include) == 0 {
		ii.Files.Include = append(ii.Files.Include, ".")
	}
	resolve := func(list []string, p string) ([]string, error) {
		if !filepath.IsAbs(p) {
			p = filepath.Join(workDir, p)
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
	for _, it := range ii.Files.Exclude {
		excluded, _ = resolve(excluded, it)
	}
	excluded = append(excluded, workOutputDir, workInstallerDir, filepath.Join(workDir, "mkinstall.json"))
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
		n = 9223372036854775807 // int64.max
	}
	return int(n), err
}

func (ii *MakeInstallInfo) InputType() string {
	buf := make([]rune, 0, 9)
	switch ii.Files.Type {
	case "raw", "tar", "zstd":
		buf = append(buf, []rune(strings.ToUpper(ii.Files.Type[0:1]))...)
		buf = append(buf, []rune(ii.Files.Type[1:])...)
	default:
		panic("unreachable")
	}
	buf = append(buf, 'I', 'n', 'p', 'u', 't')
	return string(buf)
}
