package main

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/aerogo/pixy"
	"github.com/fatih/color"
)

var scriptAnnouncePrefix = " " + color.CyanString("❄") + " "

const moduleLoader = `"use strict";

const _modules = {
${PACK_MODULES}
};

function require(path) {
	var loader = _modules[path];
	
	if(!loader)
		throw "Module not found";
	
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

func scriptFinish(results WorkerPoolResults) {
	if app.Config.Scripts.Main == "" {
		panic(errors.New("Main script file has not been defined in config.json (config.scripts.main)"))
	}

	modules := make([]string, 0, len(results))

	for job, result := range results {
		filePath := job.(string)
		code := result.(string)
		code = strings.TrimPrefix(code, `"use strict";`)
		code = strings.TrimSpace(code)
		code = strings.TrimPrefix(code, `Object.defineProperty(exports, "__esModule", { value: true });`)
		// code = strings.TrimSpace(code)

		moduleCode := `"` + strings.TrimSuffix(filePath, ".js") + `": function(exports) {` + code + "\n" + "}"
		modules = append(modules, moduleCode)

		fmt.Println(scriptAnnouncePrefix, filePath)
	}

	moduleList := strings.Join(modules, ",\n")
	bundledJS := strings.Replace(moduleLoader, "${PACK_MODULES}", moduleList, 1) + "\n" + `require("scripts/` + app.Config.Scripts.Main + `");`

	// Write JS bundle into $.js.go where it can be referenced as components.JS
	EmbedData(path.Join(outputFolder, "$.js.go"), pixy.PackageName, "JS", bundledJS)
}
