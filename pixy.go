package main

import (
	"fmt"
	"os"

	"github.com/aerogo/pixy"
	"github.com/fatih/color"
)

const outputFolder = "components"

var pixyAnnouncePrefix = " " + color.GreenString("❀") + " "

func init() {
	pixy.PackageName = outputFolder

	// Create a clean "components" directory
	os.RemoveAll(outputFolder)
	os.Mkdir(outputFolder, 0777)
}

func pixyWork(job interface{}) interface{} {
	file := job.(string)
	fmt.Println(pixyAnnouncePrefix, file)

	pixy.CompileFileAndSaveIn(file, outputFolder)

	return ""
}

func pixyFinish(results WorkerPoolResults) {
	// ...
}
