package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/OneOfOne/xxhash"
	"github.com/aerogo/pixy"
	"github.com/fatih/color"
)

const outputFolder = "components"

var pixyAnnouncePrefix = " " + color.GreenString("❀") + " "
var workDir = "./"

func init() {
	pixy.PackageName = outputFolder

	// Create a clean "components" directory
	os.RemoveAll(outputFolder)
	os.Mkdir(outputFolder, 0777)

	// Get working directory
	var err error
	workDir, err = os.Getwd()

	if err != nil {
		panic(err)
	}

	// Create pixy cache
	os.Mkdir(cacheFolder+"pixy", 0777)
}

func pixyWork(job interface{}) interface{} {
	file := job.(string)
	fmt.Println(pixyAnnouncePrefix, file)

	fullPath := path.Join(workDir, file)
	fileStat, err := os.Stat(fullPath)

	if err != nil {
		panic(err)
	}

	h := xxhash.NewS64(0)
	h.WriteString(fullPath)
	hash := strconv.FormatUint(h.Sum64(), 16)

	pixyCacheDir := path.Join(cacheFolder, "pixy", hash)

	cacheStat, cacheErr := os.Stat(pixyCacheDir)

	// Use cached version if possible
	if cacheErr == nil && cacheStat.ModTime().Unix() > fileStat.ModTime().Unix() {
		files, err := ioutil.ReadDir(pixyCacheDir)

		if err != nil {
			panic(err)
		}

		for _, file := range files {
			in := path.Join(pixyCacheDir, file.Name())
			out := path.Join(outputFolder, file.Name())

			code, err := ioutil.ReadFile(in)

			if err != nil {
				panic(err)
			}

			err = ioutil.WriteFile(out, code, 0644)

			if err != nil {
				panic(err)
			}
		}

		return ""
	}

	// We need a fresh recompile
	components := pixy.CompileFileAndSaveIn(file, outputFolder)

	// Start with an empty directory.
	// This will also reset the ModTime() of the directory.
	os.RemoveAll(pixyCacheDir)
	os.Mkdir(pixyCacheDir, 0777)

	// Cache the components in the new cache folder
	for _, component := range components {
		component.Save(pixyCacheDir)
	}

	return ""
}

func pixyFinish(results WorkerPoolResults) {
	pixy.SaveUtilities(outputFolder)
}
