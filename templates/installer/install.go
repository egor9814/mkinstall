package main

import (
	"fmt"
	"os"
	"strings"
)

type InstallInfo struct {
	ProductName        string
	TargetPath         string
	TargetPathEditable bool
	InputType          InputType
	Decrypt            bool
	Files              []string
}

func (ii *InstallInfo) init() error {
	targetPath := strings.Split(ii.TargetPath, "/")
	for i, it := range targetPath {
		if len(it) > 0 && it[0] == '$' {
			v, ok := os.LookupEnv(it[1:])
			if ok {
				targetPath[i] = v
			} else {
				return fmt.Errorf("unknown environment variable %q", it)
			}
		}
	}
	ii.TargetPath = strings.Join(targetPath, string(os.PathSeparator))

	return nil
}
