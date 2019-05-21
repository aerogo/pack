package main

import (
	"os"

	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/pack"
	jspacker "github.com/aerogo/pack/assets/js"
	pixypacker "github.com/aerogo/pack/assets/pixy"
	scarletpacker "github.com/aerogo/pack/assets/scarlet"
	"github.com/akyoto/color"
)

func main() {
	// Create a new packer
	packer := pack.New(".")

	// Initialize the asset packers
	pixy := pixypacker.New(packer.Root)
	scarlet := scarletpacker.New(packer.Root, packer.Config.Styles, packer.Config.Fonts)

	// Here we define the asset compilers.
	// Each compiler is assigned to a specific extension
	// and also has its own job queue where we will push
	// file paths as work assignments to the queue.
	packer.Compilers = []pack.AssetCompiler{
		{
			Extension:      ".pixy",
			Jobs:           jobqueue.New(pixy.Map),
			ProcessResults: pixy.Reduce,
		},
		{
			Extension:      ".scarlet",
			Jobs:           jobqueue.New(scarlet.Map),
			ProcessResults: scarlet.Reduce,
		},
	}

	// Only pack js files if the entry point has been defined
	if packer.Config.Scripts.Main != "" {
		js := jspacker.New(packer.Root, packer.Config.Scripts)

		packer.Compilers = append(packer.Compilers, pack.AssetCompiler{
			Extension:      ".js",
			Jobs:           jobqueue.New(js.Map),
			ProcessResults: js.Reduce,
		})
	}

	// They see me rollin'
	// They hatin'
	err := packer.Run()

	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
