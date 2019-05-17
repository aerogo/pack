package pack

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/akyoto/stringutils/unsafe"
)

// EmbedData embeds base64-encoded data in a Go function.
func EmbedData(outputFile, packageName, funcName, data string) error {
	outputFileAbs, err := filepath.Abs(outputFile)

	if err != nil {
		return err
	}

	// dataHash := []byte(HashString(data))
	// cacheFileName := HashString(outputFileAbs)
	// cacheFilePath := path.Join(cacheFolder, "embed", cacheFileName)

	// // Try to use cache if possible
	// cachedHash, err := ioutil.ReadFile(cacheFilePath)

	// if err == nil && bytes.Equal(dataHash, cachedHash) {
	// 	// Does the file still exist in the project directory?
	// 	if _, statErr := os.Stat(outputFileAbs); !os.IsNotExist(statErr) {
	// 		return
	// 	}
	// }

	// Encode in Base64
	base64EncodedData := base64.StdEncoding.EncodeToString(unsafe.StringToBytes(data))

	// Create Go code to load the embedded data
	loader := fmt.Sprintf(
		"package %s\n\nimport \"encoding/base64\"\n\n// %s ...\nfunc %s() string {\nencoded := `\n%s\n`\ndecoded, _ := base64.StdEncoding.DecodeString(encoded)\nreturn string(decoded)\n}\n",
		packageName,
		funcName,
		funcName,
		base64EncodedData,
	)

	// Write the cache file
	// ioutil.WriteFile(cacheFilePath, dataHash, 0644)

	// Write the loader
	return ioutil.WriteFile(outputFileAbs, unsafe.StringToBytes(loader), 0644)
}
