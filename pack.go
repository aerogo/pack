package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/aerogo/pixy"
	"github.com/fatih/color"
)

const (
	pixyExtension   = ".pixy"
	stylExtension   = ".styl"
	scriptExtension = ".ts"
	outputFolder    = "components"
	outputExtension = ".go"
)

// StylusCompileResult ...
type StylusCompileResult struct {
	file string
	css  string
}

func main() {
	pixy.PackageName = outputFolder

	os.RemoveAll(outputFolder)
	os.Mkdir(outputFolder, 0777)

	filepath.Walk(".", func(file string, f os.FileInfo, err error) error {
		if f.IsDir() || strings.HasPrefix(file, ".") {
			return nil
		}

		switch filepath.Ext(file) {
		// Template
		case pixyExtension:
			fmt.Println(" "+color.GreenString("❀"), file)
			pixy.CompileFileAndSaveIn(file, outputFolder)

		// Style
		case stylExtension:
			compileStyle(file)

			// Script
			// case scriptExtension:
			// compileScript(file)
		}

		return nil
	})

	// $.css.go
	cssCode := "package " + pixy.PackageName + "\n\nconst BundledCSS = `" + getBundledCSS() + "`\n"
	ioutil.WriteFile(path.Join(outputFolder, "$.css.go"), []byte(cssCode), 0644)

	// tscOutput, tscError := exec.Command("tsc").CombinedOutput()

	// if tscError != nil {
	// 	color.Red("Couldn't execute tsc.")
	// 	color.Red(tscError.Error())
	// }

	// fmt.Print(string(tscOutput))

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
