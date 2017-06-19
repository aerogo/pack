package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aerogo/aero"
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

// EmbedData embeds base64-encoded data in a Go function.
func EmbedData(outputFile, packageName, funcName, data string) {
	// Encode in Base64
	data = base64.StdEncoding.EncodeToString(aero.StringToBytesUnsafe(data))

	// Create Go code to load the embedded data
	loader := "package " + packageName + "\n\nimport \"encoding/base64\"\n\n// " + funcName + " ...\nfunc " + funcName + "() string {\nencoded := `\n" + data + "\n`\ndecoded, _ := base64.StdEncoding.DecodeString(encoded)\nreturn string(decoded)\n}\n"

	// Write the loader
	ioutil.WriteFile(outputFile, aero.StringToBytesUnsafe(loader), 0644)
}
