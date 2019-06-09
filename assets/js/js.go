package jspacker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/pack"
	"github.com/akyoto/color"
	"github.com/akyoto/hash"
	"github.com/akyoto/stringutils/unsafe"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

// JSPacker is a packer for javascript files.
type JSPacker struct {
	// Root directory
	root string

	// Where embedded code will be stored
	outputDirectory string

	// The prefix used for terminal output on each file.
	prefix string

	// Scripts configuration
	scripts pack.ScriptsConfiguration
}

// New creates a new JSPacker.
func New(root string, scripts pack.ScriptsConfiguration) *JSPacker {
	if scripts.Main == "" {
		panic("Main script file has not been defined in config.json (config.scripts.main)")
	}

	outputDirectory := path.Join(root, "components", "js")
	err := os.MkdirAll(outputDirectory, os.ModePerm)

	if err != nil {
		panic(err)
	}

	return &JSPacker{
		root:            root,
		outputDirectory: outputDirectory,
		prefix:          color.CyanString(" ❄ "),
		scripts:         scripts,
	}
}

// Map maps each job to its processed output.
func (packer *JSPacker) Map(job interface{}) interface{} {
	file := job.(string)
	contents, err := ioutil.ReadFile(file)

	if err != nil {
		color.Red(err.Error())
		return nil
	}

	// Skip empty files
	if len(contents) == 0 {
		return nil
	}

	// Convert it to a string
	code := unsafe.BytesToString(contents)

	// Skip files with a pack:ignore comment at the top
	if strings.HasPrefix(code, "// pack:ignore") {
		return nil
	}

	code = strings.TrimPrefix(code, `"use strict";`)
	code = strings.TrimSpace(code)
	code = strings.TrimPrefix(code, `Object.defineProperty(exports, "__esModule", { value: true });`)

	scriptDir := filepath.Dir(file)

	// Make sure we always use linux style path separators
	scriptDir = filepath.ToSlash(scriptDir)

	// TODO: This is really hacky. Replace this with a proper algorithm.
	code = strings.ReplaceAll(code, `require("./`, `require("`+scriptDir+`/`)
	code = strings.ReplaceAll(code, `require("../`, `require("`+filepath.Clean(path.Join(scriptDir, ".."))+`/`)

	return code
}

// Reduce combines all outputs.
func (packer *JSPacker) Reduce(results jobqueue.Results) {
	modules := make([]string, 0, len(results))

	for job, result := range results {
		if result == nil {
			continue
		}

		file := job.(string)
		code := result.(string)

		// Remove file extension from module path
		modulePath := strings.TrimSuffix(file, ".js")

		// Index files are implied so we don't need them in the path
		modulePath = strings.TrimSuffix(modulePath, "/index")

		// Generate module code
		moduleCode := fmt.Sprintf("\"%s\": function(exports) {%s\n}", modulePath, code)
		modules = append(modules, moduleCode)

		fmt.Println(packer.prefix, file)
	}

	// This doesn't really have any meaning besides making the order deterministic.
	// Since the order is well defined and not random, hash based caching will work.
	sort.Slice(modules, func(i, j int) bool {
		a := modules[i]
		b := modules[j]

		if len(a) == len(b) {
			return hash.String(a) < hash.String(b)
		}

		return len(a) < len(b)
	})

	moduleList := strings.Join(modules, ",\n")
	bundledJS := fmt.Sprintf("%s\nrequire(\"scripts/%s\");", strings.Replace(moduleLoader, "${PACK_MODULES}", moduleList, 1), packer.scripts.Main)

	// Minify
	m := minify.New()
	buffer := strings.Builder{}
	err := js.Minify(m, &buffer, strings.NewReader(bundledJS), nil)

	if err != nil {
		panic(err)
	}

	bundledJS = buffer.String()

	// Write JS bundle into components/js/js.go where it can be used as js.Bundle()
	embedFile := path.Join(packer.outputDirectory, "js.go")
	err = pack.EmbedData(embedFile, packer.root, "js", "Bundle", bundledJS)

	if err != nil {
		panic(err)
	}
}
