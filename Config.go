package main

import (
	"os"

	jsoniter "github.com/json-iterator/go"
)

// config is used to access the config.json data.
var config *Configuration

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

// loadConfig loads the pack configuration from the given file.
func loadConfig(fileName string) (*Configuration, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	config := &Configuration{}
	decoder := jsoniter.NewDecoder(file)
	err = decoder.Decode(config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
