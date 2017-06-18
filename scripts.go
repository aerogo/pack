package main

import (
	"fmt"

	"github.com/fatih/color"
)

var scriptAnnouncePrefix = " " + color.CyanString("❄") + " "

func scriptWork(job interface{}) interface{} {
	file := job.(string)
	scriptCode, _ := ReadFile(file)
	return scriptCode
}

func scriptFinish(results WorkerPoolResults) {
	for job, result := range results {
		fileName := job.(string)
		code := result.(string)
		code = code

		fmt.Println(scriptAnnouncePrefix, fileName)
	}
}
