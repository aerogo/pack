package pack

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/aerogo/flow/jobqueue"
	jsoniter "github.com/json-iterator/go"
)

// Packer packs the assets for your app.
type Packer struct {
	Compilers []AssetCompiler
	Root      string
	config    Configuration
}

// New creates a new packer.
func New(root string) *Packer {
	return &Packer{
		Root: root,
	}
}

// Run starts packing.
func (packer *Packer) Run() error {
	// Load configuration
	err := packer.LoadConfig(path.Join(packer.Root, "config.json"))

	if err != nil {
		return err
	}

	// Map file extensions to their corresponding worker pool
	workerPools := make(map[string]*jobqueue.JobQueue)

	for _, compiler := range packer.Compilers {
		workerPools[compiler.Extension] = compiler.Jobs
	}

	err = eachFileIn(packer.Root, func(file string) {
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

	if err != nil {
		return err
	}

	// Now that the work is queued up,
	// we can wait for each job queue to finish the work.
	for _, compiler := range packer.Compilers {
		// Wait for jobs to finish
		results := compiler.Jobs.Wait()

		// Let the compiler do compiler-specific stuff with the results
		compiler.ProcessResults(results)

		// Add an empty line separator to make the output prettier
		if len(results) > 0 {
			fmt.Println()
		}
	}

	return nil
}

// LoadConfig loads the pack configuration from the given file.
func (packer *Packer) LoadConfig(fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}

	defer file.Close()
	decoder := jsoniter.NewDecoder(file)
	return decoder.Decode(&packer.config)
}
