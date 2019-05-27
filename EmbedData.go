package pack

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/akyoto/hash"
	"github.com/akyoto/stringutils/unsafe"
)

// EmbedData embeds base64-encoded data in a Go function.
func EmbedData(outputFile, root, packageName, funcName, data string) error {
	outputFileAbs, err := filepath.Abs(outputFile)

	if err != nil {
		return err
	}

	// Try to use the cache if possible
	dataHash := hash.String(data)
	fileNameHash := hash.String(outputFileAbs)
	cacheDirectory := path.Join(root, "components", ".cache", "embed")
	err = os.MkdirAll(cacheDirectory, os.ModePerm)

	if err != nil {
		return err
	}

	cacheFilePath := path.Join(cacheDirectory, strconv.FormatInt(int64(fileNameHash), 16))
	cachedHashBytes, err := ioutil.ReadFile(cacheFilePath)

	if err == nil && binary.BigEndian.Uint64(cachedHashBytes) == dataHash {
		// Does the file still exist in the project directory?
		_, err = os.Stat(outputFileAbs)

		if !os.IsNotExist(err) {
			// Cache hit
			return err
		}
	}

	// Encode in Base64
	base64EncodedData := base64.StdEncoding.EncodeToString(unsafe.StringToBytes(data))

	// Create Go code to load the embedded data
	loader := fmt.Sprintf(
		"package %s\n\nimport \"encoding/base64\"\n\n// %s returns the bundled data.\nfunc %s() string {\n\tencoded := `%s`\n\tdecoded, _ := base64.StdEncoding.DecodeString(encoded)\n\treturn string(decoded)\n}\n",
		packageName,
		funcName,
		funcName,
		base64EncodedData,
	)

	// Write the cache file
	dataHashBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dataHashBytes, dataHash)
	err = ioutil.WriteFile(cacheFilePath, dataHashBytes, 0644)

	if err != nil {
		return err
	}

	// Write the loader
	return ioutil.WriteFile(outputFileAbs, unsafe.StringToBytes(loader), 0644)
}
