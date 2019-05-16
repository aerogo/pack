package pack

import (
	"os"

	"github.com/aerogo/flow/jobqueue"
	pixypacker "github.com/aerogo/pack/assets/pixy"
	scarletpacker "github.com/aerogo/pack/assets/scarlet"
	scriptpacker "github.com/aerogo/pack/assets/script"
	jsoniter "github.com/json-iterator/go"
)

// Packer packs the assets for your app.
type Packer struct {
	config    Configuration
	compilers []AssetCompiler
}

// New creates a new packer.
func New() *Packer {
	// Here we define the asset compilers.
	// Each compiler is assigned to a specific extension
	// and also has its own job queue where we will push
	// file paths as work assignments to the queue.
	return &Packer{
		compilers: []AssetCompiler{
			{
				Extension:      ".pixy",
				Jobs:           jobqueue.New(pixypacker.Map),
				ProcessResults: pixypacker.Reduce,
			},
			{
				Extension:      ".scarlet",
				Jobs:           jobqueue.New(scarletpacker.Map),
				ProcessResults: scarletpacker.Reduce,
			},
			{
				Extension:      ".js",
				Jobs:           jobqueue.New(scriptpacker.Map),
				ProcessResults: scriptpacker.Reduce,
			},
		},
	}
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
