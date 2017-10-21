package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/pixy"
	"github.com/fatih/color"
)

var outputFolder = "components"
var outputFolderExisted bool

var pixyAnnouncePrefix = " " + color.GreenString("❀") + " "
var workDir = "./"

func init() {
	pixy.PackageName = outputFolder

	// Create a clean "components" directory
	if _, statErr := os.Stat(outputFolder); os.IsNotExist(statErr) {
		outputFolderExisted = false

		os.Mkdir(outputFolder, 0777)
		os.Mkdir(path.Join(outputFolder, "css"), 0777)
		os.Mkdir(path.Join(outputFolder, "js"), 0777)
	} else {
		outputFolderExisted = true
	}

	// Get working directory
	var err error
	workDir, err = os.Getwd()

	if err != nil {
		panic(err)
	}

	// Create pixy cache
	os.Mkdir(path.Join(cacheFolder, "pixy"), 0777)
}

func pixyWork(job interface{}) interface{} {
	file := job.(string)
	fmt.Println(pixyAnnouncePrefix, file)

	fullPath := path.Join(workDir, file)
	fileStat, err := os.Stat(fullPath)

	if err != nil {
		panic(err)
	}

	hash := HashString(fullPath)
	pixyCacheDir := path.Join(cacheFolder, "pixy", hash)

	if outputFolderExisted {
		cacheStat, cacheErr := os.Stat(pixyCacheDir)

		// Use cached version if possible
		if cacheErr == nil && cacheStat.ModTime().Unix() > fileStat.ModTime().Unix() {
			files, err := ioutil.ReadDir(pixyCacheDir)

			if err != nil {
				panic(err)
			}

			components := []*pixy.Component{}

			for _, file := range files {
				component := strings.TrimSuffix(file.Name(), ".go")
				components = append(components, &pixy.Component{
					Name: component,
					Code: "",
				})
			}

			return components
		}
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

	return components
}

func pixyFinish(results jobqueue.Results) {
	utilitiesExist := false

	// Create a map of available components
	compiledComponents := make(map[string]bool)

	for _, result := range results {
		components := result.([]*pixy.Component)

		for _, component := range components {
			compiledComponents[component.Name] = true
		}
	}

	// Delete all components that were removed
	files, _ := ioutil.ReadDir(outputFolder)

	for _, file := range files {
		fileName := file.Name()

		if strings.HasPrefix(fileName, "$") || fileName == "css" || fileName == "js" {
			if fileName == "$.go" {
				utilitiesExist = true
			}

			continue
		}

		component := strings.TrimSuffix(fileName, ".go")
		_, exists := compiledComponents[component]

		if exists {
			continue
		}

		generatedOldFile := path.Join(outputFolder, fileName)
		os.Remove(generatedOldFile)
	}

	if utilitiesExist {
		return
	}

	pixy.SaveUtilities(outputFolder)
}
