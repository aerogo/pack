package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/pixy"
)

const (
	pixyExtension    = ".pixy"
	scarletExtension = ".scarlet"
	scriptExtension  = ".ts"
	outputFolder     = "components"
	outputExtension  = ".go"
	cacheFolder      = "/tmp/pack/"
)

var app = aero.New()
var fontsCSSChannel = make(chan string, 1)

func main() {
	pixy.PackageName = outputFolder

	// Cache folder
	os.Mkdir(cacheFolder, 0777)
	os.Mkdir(path.Join(cacheFolder, "fonts"), 0777)

	go func() {
		// FOR TESTING
		app.Config.Fonts = []string{"Ubuntu"}

		if len(app.Config.Fonts) > 0 {
			cached, err := ReadFile(path.Join(cacheFolder, "fonts", strings.Join(app.Config.Fonts, "|")+".css"))

			if err == nil {
				fontsCSSChannel <- cached
			} else {
				fontsCSSChannel <- downloadFontsCSS(app.Config.Fonts)
			}
		} else {
			fontsCSSChannel <- ""
		}
	}()

	// Output folder
	os.RemoveAll(outputFolder)
	os.Mkdir(outputFolder, 0777)

	pixyWorkerPool := NewWorkerPool(pixyWork)
	scarletWorkerPool := NewWorkerPool(scarletWork)

	scanFiles(".", func(file string) {
		switch filepath.Ext(file) {
		// Template
		case pixyExtension:
			pixyWorkerPool.Queue(file)

		// Style
		case scarletExtension:
			scarletWorkerPool.Queue(file)

		// Script
		case scriptExtension:
			// ...
		}
	})

	// Wait for all worker pools to finish
	pixyWorkerPool.Wait()
	styles := ToStringMap(scarletWorkerPool.Wait())

	// $.css.go
	bundledCSS := base64.StdEncoding.EncodeToString(aero.StringToBytesUnsafe(getBundledCSS(styles)))
	cssCode := "package " + pixy.PackageName + "\n\nimport \"encoding/base64\"\n\n// CSS ...\nfunc CSS() string {\ncssEncoded := `\n" + bundledCSS + "\n`\ncssDecoded, _ := base64.StdEncoding.DecodeString(cssEncoded)\nreturn string(cssDecoded)\n}\n"
	ioutil.WriteFile(path.Join(outputFolder, "$.css.go"), aero.StringToBytesUnsafe(cssCode), 0644)

	// Browserify & Uglify
	// cmd := exec.Command("browserify", "-o", path.Join(outputFolder, "bundle.js"), "scripts/main.js")
	// browserifyOutput, err := cmd.CombinedOutput()
	// fmt.Print(string(browserifyOutput))

	// if err != nil {
	// 	color.Red("Couldn't execute browserify.")
	// 	color.Red(err.Error())
	// }

	fmt.Println()
	fmt.Println("Done.")
}
