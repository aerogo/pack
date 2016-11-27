package main

import (
	"fmt"

	"github.com/aerogo/pixy"
	"github.com/fatih/color"
)

func pixyWork(job interface{}) interface{} {
	file := job.(string)
	fmt.Println(" "+color.GreenString("❀"), file)

	pixy.CompileFileAndSaveIn(file, outputFolder)

	return "done"
}
