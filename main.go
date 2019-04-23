package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/aerogo/aero"
	"github.com/aerogo/flow/jobqueue"
)

// Cache folder is used as our location for cache files.
var cacheFolder = path.Join(os.TempDir(), "pack")

// config is used to access the aero config.json data.
var config *aero.Configuration

func main() {
	// Load config file
	var err error
	config, err = aero.LoadConfig("config.json")
	PanicOnError(err)

	// Cache folder
	os.Mkdir(cacheFolder, 0777)

	// Init
	fontsInit()
	pixyInit()
	embedInit()

	// Here we define the asset compilers.
	// Each compiler is assigned to a specific extension
	// and also has its own job queue where we will push
	// file paths as work assignments to the queue.
	compilers := []AssetCompiler{
		{
			Extension:      ".pixy",
			Jobs:           jobqueue.New(pixyWork),
			ProcessResults: pixyFinish,
		},
		{
			Extension:      ".scarlet",
			Jobs:           jobqueue.New(scarletWork),
			ProcessResults: scarletFinish,
		},
		{
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

	// Scan the current directory and assign work to each file by file extension
	ScanFiles(".", func(file string) {
		// Check if we have a compiler registered for that file type
		workerPool, exists := workerPools[filepath.Ext(file)]

		if !exists {
			return
		}

		// Make sure we always use linux style path separators
		file = filepath.ToSlash(file)

		// Queue up work by sending the file path to the compiler
		workerPool.Queue(file)
	})

	// Now that the work is queued up,
	// we can wait for each job queue to finish the work.
	for _, compiler := range compilers {
		// Wait for jobs to finish
		results := compiler.Jobs.Wait()

		// Let the compiler do compiler-specific things with the results
		compiler.ProcessResults(results)

		// Add an empty line separator to make the output prettier
		if len(results) > 0 {
			println()
		}
	}
}
