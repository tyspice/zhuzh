package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type ChatGPT struct {
	ApiKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

type Config struct {
	ChatGPT ChatGPT `yaml:"chatgpt"`
}

var (
	config Config
	mu     sync.RWMutex
)

func init() {
	config = Config{}

	// Add defaults
	config.ChatGPT.Model = "gpt-3.5-turbo"

	// Read in yaml config
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(homeDir, ".config", "zhuzh", "config.yml")
		if data, err := os.ReadFile(configPath); err == nil {
			// Ignore yaml parsing errors and continue with defaults
			yaml.Unmarshal(data, &config)
		}
	}
}

func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

func GetChatGPT() ChatGPT {
	mu.RLock()
	defer mu.RUnlock()
	return config.ChatGPT
}

func Set(c Config) {
	mu.Lock()
	defer mu.Unlock()
	config = c
}
