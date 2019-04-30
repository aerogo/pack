package main

import (
	"bufio"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aerogo/flow/jobqueue"
	"github.com/akyoto/color"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

var scriptAnnouncePrefix = " " + color.CyanString("❄") + " "

const moduleLoader = `"use strict";

const _modules = {
${PACK_MODULES}
};

function require(path) {
	var loader = _modules[path];
	
	if(!loader)
		throw "Module not found: " + path;
	
	if(loader.exports !== undefined)
		return loader.exports;
	
	loader.exports = {};
	loader(loader.exports);
	
	return loader.exports;
}`

func scriptWork(job interface{}) interface{} {
	file := job.(string)
	scriptCode, _ := ReadFile(file)
	return scriptCode
}

func scriptFinish(results jobqueue.Results) {
	if config.Scripts.Main == "" {
		panic("Main script file has not been defined in config.json (config.scripts.main)")
	}

	modules := make([]string, 0, len(results))

	for job, result := range results {
		file := job.(string)
		code := result.(string)

		// Skip empty files
		if len(code) == 0 {
			continue
		}

		// Module that have the pack:ignore comment at the top will be ignored
		if strings.HasPrefix(code, "// pack:ignore") {
			continue
		}

		code = strings.TrimPrefix(code, `"use strict";`)
		code = strings.TrimSpace(code)
		code = strings.TrimPrefix(code, `Object.defineProperty(exports, "__esModule", { value: true });`)
		// code = strings.TrimSpace(code)

		scriptDir := filepath.Dir(file)

		// Normalize file paths (Windows)
		scriptDir = strings.Replace(scriptDir, "\\", "/", -1)

		// TODO: This is really hacky. Replace this with a proper algorithm.
		code = strings.Replace(code, `require("./`, `require("`+scriptDir+`/`, -1)
		code = strings.Replace(code, `require("../`, `require("`+filepath.Clean(path.Join(scriptDir, ".."))+`/`, -1)

		// Remove file extension from module path
		modulePath := strings.TrimSuffix(file, ".js")

		// Index files are implied so we don't need them in the path
		modulePath = strings.TrimSuffix(modulePath, "/index")

		moduleCode := `"` + modulePath + `": function(exports) {` + code + "\n" + "}"
		modules = append(modules, moduleCode)

		fmt.Println(scriptAnnouncePrefix, file)
	}

	// This doesn't really have any meaning besides making the order deterministic.
	// Since the order is well defined and not random, hash based caching will work.
	sort.Slice(modules, func(i, j int) bool {
		a := modules[i]
		b := modules[j]

		if len(a) == len(b) {
			return HashString(a) < HashString(b)
		}

		return len(a) < len(b)
	})

	moduleList := strings.Join(modules, ",\n")
	bundledJS := strings.Replace(moduleLoader, "${PACK_MODULES}", moduleList, 1) + "\n" + `require("scripts/` + config.Scripts.Main + `");`

	// // Minify
	m := minify.New()
	buffer := strings.Builder{}
	writer := bufio.NewWriter(&buffer)
	err := js.Minify(m, writer, strings.NewReader(bundledJS), nil)

	if err != nil {
		panic(err)
	}

	writer.Flush()
	bundledJS = buffer.String()

	// Write JS bundle into $.js.go where it can be referenced as components.JS
	EmbedData(path.Join(outputFolder, "js", "js.go"), "js", "Bundle", bundledJS)
}
