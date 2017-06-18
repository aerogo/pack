package main

import (
	"path/filepath"

	"github.com/aerogo/aero"
)

const cacheFolder = "/tmp/pack/"

var app = aero.New()

func main() {
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
}
