package scarletpacker

import (
	"fmt"

	"github.com/aerogo/flow/jobqueue"
	"github.com/akyoto/color"
)

// ScarletPacker is a packer for scarlet files.
type ScarletPacker struct {
	// Root directory
	root string

	// The prefix used for terminal output on each file.
	prefix string
}

// New creates a new ScarletPacker.
func New(root string) *ScarletPacker {
	return &ScarletPacker{
		root:   root,
		prefix: color.YellowString(" ★ "),
	}
}

// Map maps each job to its processed output.
func (packer *ScarletPacker) Map(job interface{}) interface{} {
	// fileName := job.(string)
	// fmt.Println(packer.prefix, fileName)
	return nil
}

// Reduce combines all outputs.
func (packer *ScarletPacker) Reduce(results jobqueue.Results) {
	// Unordered styles in styles directory
	for styleName := range results {
		fmt.Println(packer.prefix, styleName)
	}
}
