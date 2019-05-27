package pixypacker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/pixy"
	"github.com/akyoto/autoimport"
	"github.com/akyoto/color"
	"github.com/akyoto/hash"
	"github.com/akyoto/stringutils/unsafe"
)

// PixyPacker is a packer for pixy files.
type PixyPacker struct {
	// Root directory
	root string

	// Where generated code will be stored
	outputDirectory string

	// The compiler for pixy files which
	// is initialized with the package name.
	compiler *pixy.Compiler

	// Auto importer for pixy code
	importer *autoimport.AutoImport

	// The prefix used for terminal output on each file.
	prefix string
}

// New creates a new PixyPacker.
func New(root string) *PixyPacker {
	outputDirectory := path.Join(root, "components")
	err := os.MkdirAll(outputDirectory, os.ModePerm)

	if err != nil {
		panic(err)
	}

	rootAbs, err := filepath.Abs(root)

	if err != nil {
		panic(err)
	}

	return &PixyPacker{
		root:            root,
		outputDirectory: outputDirectory,
		compiler:        pixy.NewCompiler("components"),
		importer:        autoimport.New(rootAbs),
		prefix:          color.GreenString(" ✿ "),
	}
}

// Map maps each job to its processed output.
func (packer *PixyPacker) Map(job interface{}) interface{} {
	fileName := job.(string)
	fmt.Println(packer.prefix, fileName)

	fileStat, err := os.Stat(fileName)

	if err != nil {
		panic(err)
	}

	// Each file is cached in a directory that saves
	// each component inside the file as Go code.
	fileNameHash := hash.String(fileName)
	cacheDirectory := path.Join(packer.root, "components", ".cache", "pixy", strconv.FormatInt(int64(fileNameHash), 16))
	cacheStat, cacheErr := os.Stat(cacheDirectory)

	// Use cached version if possible
	if cacheErr == nil && cacheStat.ModTime().Unix() > fileStat.ModTime().Unix() {
		cachedComponentFiles, err := ioutil.ReadDir(cacheDirectory)

		if err != nil {
			panic(err)
		}

		// We will fake a list of compiled components by searching
		// the cache directory and returning each file name without
		// the ".go" suffix as a component name.
		components := make([]*pixy.Component, 0, len(cachedComponentFiles))

		for _, file := range cachedComponentFiles {
			componentName := strings.TrimSuffix(file.Name(), ".go")

			components = append(components, &pixy.Component{
				Name: componentName,
			})
		}

		return components
	}

	// We need a fresh recompile
	components, _, err := packer.compileFileAndSaveIn(fileName, packer.outputDirectory)

	if err != nil {
		color.Red(err.Error())
		return nil
	}

	// Now we need to save the compilation result to disk.
	// First off, start with an empty directory.
	// This will also reset the ModTime() of the directory.
	err = os.RemoveAll(cacheDirectory)

	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(cacheDirectory, os.ModePerm)

	if err != nil {
		panic(err)
	}

	// Cache the components in the new cache folder
	for _, component := range components {
		packer.saveComponent(component, cacheDirectory)
	}

	return components
}

// Reduce combines all outputs.
func (packer *PixyPacker) Reduce(results jobqueue.Results) {
	utilsFile := path.Join(packer.outputDirectory, "utils.go")
	err := packer.compiler.SaveUtilities(utilsFile)

	if err != nil {
		panic(err)
	}
}

// compileFileAndSaveIn compiles a pixy template from fileIn
// and writes the resulting components to dirOut.
func (packer *PixyPacker) compileFileAndSaveIn(fileIn string, dirOut string) ([]*pixy.Component, []string, error) {
	components, err := packer.compiler.CompileFile(fileIn)
	files := make([]string, len(components))

	for index, component := range components {
		files[index] = packer.saveComponent(component, dirOut)
	}

	return components, files, err
}

// saveComponent writes the component to the given
// directory and returns the file path.
func (packer *PixyPacker) saveComponent(component *pixy.Component, dirOut string) string {
	file := path.Join(dirOut, component.Name+".go")
	newCode, err := packer.importer.Source(unsafe.StringToBytes(component.Code))

	if err != nil {
		color.Red("Can't autoimport " + file)
		color.Red(err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(file, newCode, 0644)

	if err != nil {
		color.Red("Can't write to " + file)
		color.Red(err.Error())
		os.Exit(1)
	}

	return file
}
