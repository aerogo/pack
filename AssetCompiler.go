package main

// AssetCompiler ...
type AssetCompiler struct {
	Extension      string
	WorkerPool     *WorkerPool
	ProcessResults func(WorkerPoolResults)
}
