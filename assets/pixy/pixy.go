package pixypacker

import (
	"fmt"

	"github.com/aerogo/flow/jobqueue"
	"github.com/aerogo/pixy"
	"github.com/akyoto/color"
)

// PixyPacker is a packer for pixy files.
type PixyPacker struct {
	// Root directory
	root string

	// The prefix used for terminal output on each file.
	prefix string
}

// New creates a new PixyPacker.
func New(root string) *PixyPacker {
	return &PixyPacker{
		root:   root,
		prefix: color.GreenString(" ✿ "),
	}
}

// Map maps each job to its processed output.
func (packer *PixyPacker) Map(job interface{}) interface{} {
	fileName := job.(string)
	fmt.Println(packer.prefix, fileName)
	components, err := pixy.CompileFile(fileName)

	if err != nil {
		color.Red(err.Error())
		return nil
	}

	return components
}

// Reduce combines all outputs.
func (packer *PixyPacker) Reduce(results jobqueue.Results) {
}
