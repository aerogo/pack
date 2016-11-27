package main

import (
	"fmt"

	"github.com/aerogo/pixy"
	"github.com/fatih/color"
)

var pixyAnnouncePrefix = " " + color.GreenString("❀") + " "

func pixyWork(job interface{}) interface{} {
	file := job.(string)
	fmt.Println(pixyAnnouncePrefix, file)

	pixy.CompileFileAndSaveIn(file, outputFolder)

	return "done"
}
