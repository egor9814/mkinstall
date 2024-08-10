//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func _create_shortcut(s *Shortcut) error {
	// https://stackoverflow.com/a/32467891
	script := fmt.Sprintf(`
option explicit

sub CreateShortCut()
	dim objShell, strDesktopPath, objLink
	set objShell = CreateObject("WScript.Shell")
	strDesktopPath = objShell.SpecialFolders("Desktop")
	set objLink = objShell.CreateShortcut(strDesktopPath & "\%s.lnk")
	objLink.Arguments = "%s"
	objLink.TargetPath = "%s"
	objLink.WindowStyle = 1
	objLink.WorkingDirectory = "%s"
	objLink.Save
end sub

call CreateShortCut()
`, s.Name, strings.Join(s.Arguments, " "), filepath.Join(install.TargetPath, filepath.FromSlash(s.Target)), install.TargetPath)

	var scriptFile string
	for {
		scriptFile = filepath.Join(install.TargetPath, "__temp_"+strconv.Itoa(int(time.Now().Unix()))+".vbs")
		if _, err := os.Stat(scriptFile); err != nil {
			break
		}
	}
	defer func() {
		os.Remove(scriptFile)
	}()
	if err := os.WriteFile(scriptFile, []byte(script), 0700); err != nil {
		return err
	}

	cmd := exec.Command("wscript", "//Nologo", scriptFile)
	return cmd.Run()
}
