package main

import (
	"os"
	"path/filepath"
	"strings"
)

func scanFiles(dir string, cb func(string)) {
	filepath.Walk(dir, func(file string, f os.FileInfo, err error) error {
		if f.IsDir() || strings.HasPrefix(file, ".") {
			return nil
		}

		cb(file)

		return nil
	})
}
