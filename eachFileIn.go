package pack

import (
	"os"
	"path/filepath"
	"strings"
)

// eachFileIn calls the callback function on each file in the given directory.
func eachFileIn(dir string, onFile func(string)) error {
	return filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if info.IsDir() || strings.HasPrefix(file, ".") {
			return nil
		}

		onFile(file)
		return nil
	})
}
