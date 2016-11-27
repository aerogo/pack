package main

import (
	"os/exec"

	"github.com/fatih/color"
)

func compileScript(file string) {
	// _, err := exec.Command("tsc", "--target", "ES6", "--outFile", path.Join(".tmp", filepath.(file)), file).Output()
	_, err := exec.Command("tsc", "--outDir", ".tmp", file).Output()

	if err != nil {
		color.Red("Couldn't execute tsc.")
		color.Red(err.Error())
		return
	}

	// tsc --target "ES6" --outDir .scripts/test --baseUrl scripts scripts/posts.ts
	// browserify .scripts/*.js -o bundle.js
	// uglifyjs --screw-ie8
}

// Browserify & Uglify
// cmd := exec.Command("browserify", "-o", path.Join(outputFolder, "bundle.js"), "scripts/main.js")
// browserifyOutput, err := cmd.CombinedOutput()
// fmt.Print(string(browserifyOutput))

// if err != nil {
// 	color.Red("Couldn't execute browserify.")
// 	color.Red(err.Error())
// }
