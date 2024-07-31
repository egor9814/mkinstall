package main

import (
	"fmt"
	"os"
)

func cleanup() error {
	errorList := make([]error, 0, 2)
	for _, it := range []string{workInstallerDir, workOutputDir} {
		if err := os.RemoveAll(it); err != nil {
			errorList = append(errorList, err)
		}
	}
	if l := len(errorList); l == 0 {
		return nil
	} else if l == 1 {
		return errorList[0]
	} else {
		return fmt.Errorf("%v\n%v", errorList[0], errorList[1])
	}
}
