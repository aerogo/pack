package main

import (
	"os"
	"path/filepath"

	"github.com/aerogo/aero"
	"github.com/aerogo/pixy"
)

const (
	cacheFolder      = "/tmp/pack/"
	outputFolder     = "components"
	outputExtension  = ".go"
	scriptExtension  = ".js"
	pixyExtension    = ".pixy"
	scarletExtension = ".scarlet"
)

var app = aero.New()

func main() {
	pixy.PackageName = outputFolder

	// Create a clean "components" directory
	os.RemoveAll(outputFolder)
	os.Mkdir(outputFolder, 0777)

	// Compilers
	compilers := []AssetCompiler{
		AssetCompiler{
			Extension:      ".pixy",
			WorkerPool:     NewWorkerPool(pixyWork),
			ProcessResults: pixyFinish,
		},
		AssetCompiler{
			Extension:      ".scarlet",
			WorkerPool:     NewWorkerPool(scarletWork),
			ProcessResults: scarletFinish,
		},
		AssetCompiler{
			Extension:      ".js",
			WorkerPool:     NewWorkerPool(scriptWork),
			ProcessResults: scriptFinish,
		},
	}

	// Map file extensions to their corresponding worker pool
	workerPools := make(map[string]*WorkerPool)

	for _, compiler := range compilers {
		workerPools[compiler.Extension] = compiler.WorkerPool
	}

	// Assign work by file extension
	scanFiles(".", func(file string) {
		workerPool, exists := workerPools[filepath.Ext(file)]

		if !exists {
			return
		}

		workerPool.Queue(file)
	})

	for _, compiler := range compilers {
		results := compiler.WorkerPool.Wait()
		compiler.ProcessResults(results)
		println()
	}

	// // Wait for all pixy workers to finish
	// pixyWorkerPool.Wait()
	// fmt.Println("")

	// // Scripts
	// scripts := ToStringMap(scriptWorkerPool.Wait())
	// for name := range scripts {
	// 	fmt.Println(scriptAnnouncePrefix, name)
	// }

	// fmt.Println()
	// fmt.Println("Done.")
}
