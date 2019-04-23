package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/akyoto/autoimport"

	"github.com/aerogo/aero"
	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/pixy"
	"github.com/akyoto/color"
)

var (
	// Output folder is the name of the generated components directory.
	outputFolder = "components"

	// This flag tells us whether the output folder existed at runtime.
	// If it did not exist, then we don't need to check for a cached version.
	outputFolderExisted = false

	// The compiler for pixy files is initialized with the package name.
	pixyCompiler = pixy.NewCompiler(outputFolder)

	// Auto importer for pixy code
	pixyImporter *autoimport.AutoImport

	// Pixy announce prefix is the prefix used for terminal output on each file.
	pixyAnnouncePrefix = " " + color.GreenString("❀") + " "

	// This will contain the current working directory.
	workDir = "./"
)

type pixyCompilationResult struct {
	Components []*pixy.Component
	Files      []string
}

func pixyInit() {
	// Create a clean "components" directory
	if _, statErr := os.Stat(outputFolder); os.IsNotExist(statErr) {
		os.Mkdir(outputFolder, 0777)
		os.Mkdir(path.Join(outputFolder, "css"), 0777)
		os.Mkdir(path.Join(outputFolder, "js"), 0777)
	} else {
		outputFolderExisted = true
	}

	// Get working directory
	var err error
	workDir, err = os.Getwd()
	PanicOnError(err)

	// Create importer
	pixyImporter = autoimport.New(workDir)

	// Create pixy cache
	os.Mkdir(path.Join(cacheFolder, "pixy"), 0777)
}

func pixyWork(job interface{}) interface{} {
	file := job.(string)
	fmt.Println(pixyAnnouncePrefix, file)

	fullPath := path.Join(workDir, file)
	fileStat, err := os.Stat(fullPath)
	PanicOnError(err)

	hash := HashString(fullPath)
	pixyCacheDir := path.Join(cacheFolder, "pixy", hash)

	if outputFolderExisted {
		cacheStat, cacheErr := os.Stat(pixyCacheDir)

		// Use cached version if possible
		if cacheErr == nil && cacheStat.ModTime().Unix() > fileStat.ModTime().Unix() {
			files, err := ioutil.ReadDir(pixyCacheDir)
			PanicOnError(err)

			components := []*pixy.Component{}

			for _, file := range files {
				component := strings.TrimSuffix(file.Name(), ".go")
				components = append(components, &pixy.Component{
					Name: component,
				})
			}

			return &pixyCompilationResult{
				Components: components,
			}
		}
	}

	// We need a fresh recompile
	components, files, err := compileFileAndSaveIn(pixyCompiler, file, outputFolder)

	if err != nil {
		color.Red(err.Error())
		return nil
	}

	// Start with an empty directory.
	// This will also reset the ModTime() of the directory.
	os.RemoveAll(pixyCacheDir)
	os.Mkdir(pixyCacheDir, 0777)

	// Cache the components in the new cache folder
	for _, component := range components {
		savePixyComponent(component, pixyCacheDir)
	}

	return &pixyCompilationResult{
		Components: components,
		Files:      files,
	}
}

func pixyFinish(results jobqueue.Results) {
	utilitiesExist := false

	// Create a map of available components
	compiledComponents := make(map[string]bool)
	var writtenFiles []string

	for _, obj := range results {
		result := obj.(*pixyCompilationResult)

		for _, component := range result.Components {
			compiledComponents[component.Name] = true
		}

		writtenFiles = append(writtenFiles, result.Files...)
	}

	// Delete all components that were removed
	files, _ := ioutil.ReadDir(outputFolder)

	for _, file := range files {
		fileName := file.Name()

		if fileName == "utils.go" {
			utilitiesExist = true
			continue
		}

		if !strings.HasSuffix(fileName, ".go") {
			continue
		}

		component := strings.TrimSuffix(fileName, ".go")
		_, exists := compiledComponents[component]

		if exists {
			continue
		}

		generatedOldFile := path.Join(outputFolder, fileName)
		err := os.Remove(generatedOldFile)

		if err != nil {
			color.Red(err.Error())
		}
	}

	if utilitiesExist {
		return
	}

	pixyCompiler.SaveUtilities(path.Join(outputFolder, "utils.go"))
}

// compileFileAndSaveIn compiles a pixy template from fileIn
// and writes the resulting components to dirOut.
func compileFileAndSaveIn(compiler *pixy.Compiler, fileIn string, dirOut string) ([]*pixy.Component, []string, error) {
	components, err := compiler.CompileFile(fileIn)
	files := make([]string, len(components))

	for index, component := range components {
		files[index] = savePixyComponent(component, dirOut)
	}

	return components, files, err
}

// savePixyComponent writes the component to the given directory and returns the file path.
func savePixyComponent(component *pixy.Component, dirOut string) string {
	file := path.Join(dirOut, component.Name+".go")
	newCode, err := pixyImporter.Source(aero.StringToBytesUnsafe(component.Code))

	if err != nil {
		color.Red("Can't autoimport " + file)
		color.Red(err.Error())
	}

	err = ioutil.WriteFile(file, newCode, 0644)

	if err != nil {
		color.Red("Can't write to " + file)
		color.Red(err.Error())
	}

	return file
}
