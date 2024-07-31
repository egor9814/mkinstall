package main

import (
	_ "embed"
	"log"
	"strconv"
	"strings"
)

//go:generate go run tools/get_version.go

//go:embed version.ref
var version_ref string

var Version struct {
	Major, Minor, Patch int
	Suffix              string
}

func parseVersion() {
	s := version_ref
	if len(s) == 0 {
		log.Fatal("version.ref is empty")
	}
	if s[0] == 'v' {
		s = s[1:]
	}
	i := strings.Index(s, "-")
	j := strings.Index(s, "+")
	if i == -1 {
		i = j
	} else if j != -1 {
		i = min(i, j)
	}
	if i != -1 {
		Version.Suffix = s[i:]
		s = s[:i]
	}
	numbers := strings.Split(s, ".")
	output := []*int{
		&Version.Major, &Version.Minor, &Version.Patch,
	}
	for i, it := range numbers {
		if n, err := strconv.ParseInt(it, 10, 64); err != nil {
			log.Fatal(err)
		} else {
			*output[i] = int(n)
		}
	}
}
