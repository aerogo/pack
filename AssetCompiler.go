package pack

import (
	"github.com/aerogo/flow/jobqueue"
)

// AssetCompiler represents a compiler for a group of assets.
type AssetCompiler struct {
	Extension      string
	Jobs           *jobqueue.JobQueue
	ProcessResults func(jobqueue.Results)
}
