//go:build !windows

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func _get_home_dir() (home string, err error) {
	var ok bool
	home, ok = os.LookupEnv("HOME")
	if !ok {
		err = errors.New("user home dir not specified")
	}
	return
}

func _get_dekstop_from_user_dirs() (string, error) {
	home, err := _get_home_dir()
	if err != nil {
		return "", err
	}
	// https://unix.stackexchange.com/a/545228
	file := filepath.Join(home, ".config", "user-dirs.dirs")
	if _, err := os.Stat(file); err != nil {
		return "", os.ErrNotExist
	}
	contentBytes, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	content := []rune(string(contentBytes))
	var begin int
	for i, l := 0, len(content); i < l; i++ {
		for ; i < l && content[i] != '=' && content[i] != '\n'; i++ {
		}
		if i == l {
			break
		}
		if content[i] == '\n' {
			begin = i + 1
			continue
		}
		if content[i] == '=' {
			if string(content[begin:i]) == "XDG_DESKTOP_DIR" {
				i++
				begin = i
				for ; i < l && content[i] != '\n'; i++ {
				}
				if i == l || content[i] == '\n' {
					return string(content[begin:i]), nil
				}
			} else {
				continue
			}
		}
	}
	return "", os.ErrNotExist
}

func _get_desktop_dir() (string, error) {
	home, err := _get_home_dir()
	if err != nil {
		return "", err
	}
	desktopDir := filepath.Join(home, "Desktop")
	if _, err := os.Stat(desktopDir); err == nil {
		return desktopDir, nil
	}
	desktopDir2, err := _get_dekstop_from_user_dirs()
	if err == nil {
		desktopDir2 = strings.Replace(desktopDir2, "$HOME", home, 1)
		if _, err := os.Stat(desktopDir2); err == nil {
			return desktopDir2, nil
		}
	}
	if err := os.MkdirAll(desktopDir, 0744); err == nil {
		return desktopDir, nil
	}
	return "", os.ErrNotExist
}

func _create_shortcut(s *Shortcut) error {
	desktop := fmt.Sprintf(`Type=Application
Name=%q
Icon=%q
Exec="%s %s"
Path=%q
Categories=%q
`, s.Name, s.Icon, s.Target, strings.Join(s.Arguments, " "), install.TargetPath, strings.Join(s.Categories, ","))
	desktopDir, err := _get_desktop_dir()
	if err != nil {
		return err
	}
	desktopFile := filepath.Join(desktopDir, s.Name+".desktop")
	return os.WriteFile(desktopFile, []byte(desktop), 0644)
}
