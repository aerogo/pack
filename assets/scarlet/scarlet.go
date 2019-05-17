package scarletpacker

import (
	"fmt"

	"github.com/aerogo/flow/jobqueue"
	"github.com/akyoto/color"
)

var (
	// The prefix used for terminal output on each file.
	prefix = " " + color.YellowString("★") + " "
)

func Map(job interface{}) interface{} {
	fmt.Println(prefix, job)
	return nil
}

func Reduce(results jobqueue.Results) {
}
