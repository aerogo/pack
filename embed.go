package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/aerogo/aero"
)

func init() {
	// Create embed cache
	os.Mkdir(path.Join(cacheFolder, "embed"), 0777)
}

// EmbedData embeds base64-encoded data in a Go function.
func EmbedData(outputFile, packageName, funcName, data string) {
	dataHash := []byte(HashString(data))
	outputFileAbs, _ := filepath.Abs(outputFile)
	cacheFileName := HashString(outputFileAbs)
	cacheFilePath := path.Join(cacheFolder, "embed", cacheFileName)

	// Try to use cache if possible
	cachedHash, err := ioutil.ReadFile(cacheFilePath)

	if err == nil && bytes.Equal(dataHash, cachedHash) {
		// Does the file still exist in the project directory?
		if _, statErr := os.Stat(outputFileAbs); !os.IsNotExist(statErr) {
			return
		}
	}

	// Encode in Base64
	data = base64.StdEncoding.EncodeToString(aero.StringToBytesUnsafe(data))

	// Create Go code to load the embedded data
	loader := "package " + packageName + "\n\nimport \"encoding/base64\"\n\n// " + funcName + " ...\nfunc " + funcName + "() string {\nencoded := `\n" + data + "\n`\ndecoded, _ := base64.StdEncoding.DecodeString(encoded)\nreturn string(decoded)\n}\n"

	// Write the loader
	ioutil.WriteFile(outputFileAbs, aero.StringToBytesUnsafe(loader), 0644)

	// Write the cache file
	ioutil.WriteFile(cacheFilePath, dataHash, 0644)
}
