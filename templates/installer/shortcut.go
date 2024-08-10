package main

import "fmt"

type Shortcut struct {
	Name       string
	Target     string
	Arguments  []string
	Icon       string
	Categories []string
}

func (s *Shortcut) Make() (err error) {
	err = _create_shortcut(s)
	if err != nil {
		err = fmt.Errorf("cannot create shortcut file: %v", err)
	}
	return
}
