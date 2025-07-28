package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ChatGPT struct {
		ApiKey string `yaml:"api_key"`
		Model  string `yaml:"model"`
	} `yaml:"chatgpt"`
}

type ConfigOption struct {
	ChatGPTApiKey *string
	ChatGPTModel  *string
}

var (
	instance *Config
	once     sync.Once
)

// Init creates the singleton instance by first reading from ~/.config/zhuzh/config.yml
// and then overriding with any provided options
// This should be called once at application startup
func Init(opts *ConfigOption) *Config {
	once.Do(func() {
		instance = &Config{}

		// Read from YAML file first
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configPath := filepath.Join(homeDir, ".config", "zhuzh", "config.yml")
			if data, err := os.ReadFile(configPath); err == nil {
				// Ignore YAML parsing errors and continue with defaults
				yaml.Unmarshal(data, instance)
			}
		}

		// Override with provided options
		if opts != nil {
			if opts.ChatGPTApiKey != nil {
				instance.ChatGPT.ApiKey = *opts.ChatGPTApiKey
			}
			if opts.ChatGPTModel != nil {
				instance.ChatGPT.Model = *opts.ChatGPTModel
			}
		}

	})
	return instance
}

// Get retrieves the singleton instance. If the singleton
// has not been initialized then we initialize it without options.
func Get() *Config {
	if instance == nil {
		return Init(nil)
	}
	return instance
}
