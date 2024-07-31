package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//go:embed install.json
var installJson string

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

func (ii *InstallInfo) load() error {
	if err := json.Unmarshal([]byte(installJson), ii); err != nil {
		return err
	}

	switch ii.Files.Type {
	case "raw", "tar", "zstd":

	default:
		return fmt.Errorf("unsupported file type %q", ii.Files.Type)
	}

	targetPath := strings.Split(ii.Target.Path, "/")
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
	ii.Target.Path = strings.Join(targetPath, string(os.PathSeparator))

	return nil
}
