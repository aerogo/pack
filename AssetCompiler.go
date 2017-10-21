package main

import (
	"github.com/aerogo/flow/jobqueue"
)

// AssetCompiler ...
type AssetCompiler struct {
	Extension      string
	Jobs           *jobqueue.JobQueue
	ProcessResults func(jobqueue.Results)
}
