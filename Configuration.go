package pack

// Configuration represents the configuration for pack.
type Configuration struct {
	Fonts   []string             `json:"fonts"`
	Styles  []string             `json:"styles"`
	Scripts ScriptsConfiguration `json:"scripts"`
}

// ScriptsConfiguration lets you configure your main entry script.
type ScriptsConfiguration struct {
	// Entry point for scripts
	Main string `json:"main"`
}
