package jspacker

import (
	"fmt"

	"github.com/aerogo/flow/jobqueue"
	"github.com/akyoto/color"
)

// JSPacker is a packer for javascript files.
type JSPacker struct {
	// Root root
	root string

	// The prefix used for terminal output on each file.
	prefix string
}

// New creates a new JSPacker.
func New(root string) *JSPacker {
	return &JSPacker{
		root:   root,
		prefix: color.CyanString(" ❄ "),
	}
}

// Map maps each job to its processed output.
func (packer *JSPacker) Map(job interface{}) interface{} {
	// fileName := job.(string)
	// fmt.Println(packer.prefix, fileName)
	return nil
}

// Reduce combines all outputs.
func (packer *JSPacker) Reduce(results jobqueue.Results) {
	// Unordered styles in styles directory
	for styleName := range results {
		fmt.Println(packer.prefix, styleName)
	}
}
