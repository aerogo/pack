package main

import (
	"path/filepath"

	"github.com/aerogo/aero"
	"github.com/aerogo/flow/jobqueue"
)

const cacheFolder = "/tmp/pack/"

var app = aero.New()

func main() {
	// Compilers
	compilers := []AssetCompiler{
		AssetCompiler{
			Extension:      ".pixy",
			Jobs:           jobqueue.New(pixyWork),
			ProcessResults: pixyFinish,
		},
		AssetCompiler{
			Extension:      ".scarlet",
			Jobs:           jobqueue.New(scarletWork),
			ProcessResults: scarletFinish,
		},
		AssetCompiler{
			Extension:      ".js",
			Jobs:           jobqueue.New(scriptWork),
			ProcessResults: scriptFinish,
		},
	}

	// Map file extensions to their corresponding worker pool
	workerPools := make(map[string]*jobqueue.JobQueue)

	for _, compiler := range compilers {
		workerPools[compiler.Extension] = compiler.Jobs
	}

	// Assign work by file extension
	ScanFiles(".", func(file string) {
		workerPool, exists := workerPools[filepath.Ext(file)]

		if !exists {
			return
		}

		workerPool.Queue(file)
	})

	for _, compiler := range compilers {
		results := compiler.Jobs.Wait()
		compiler.ProcessResults(results)
		println()
	}
}
