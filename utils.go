package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ReadFile reads in a file as a string
func ReadFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	return string(b), err
}

// ToStringMap converts map[interface{}]interface{} to map[string]string.
func ToStringMap(old map[interface{}]interface{}) map[string]string {
	newMap := make(map[string]string)

	for k, v := range old {
		newMap[k.(string)] = v.(string)
	}

	return newMap
}

// ScanFiles calls the callback function on each file in the given directory.
func ScanFiles(dir string, cb func(string)) {
	filepath.Walk(dir, func(file string, f os.FileInfo, err error) error {
		if f.IsDir() || strings.HasPrefix(file, ".") {
			return nil
		}

		cb(file)

		return nil
	})
}