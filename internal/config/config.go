package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	EnvConfigFile = "TFMGR_CONFIG"
	EnvPath       = "TFMGR_PATH"
)

type Config struct {
	Path           string `json:"path,omitempty"`
	DefaultVersion string `json:"default_version,omitempty"`
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".tfmgr")
}

func Path() string {
	if path, hasPath := os.LookupEnv(EnvPath); hasPath {
		return path
	}

	return DefaultPath()
}

func ConfigFile() string {
	if path, hasPath := os.LookupEnv(EnvConfigFile); hasPath {
		return path
	}

	return filepath.Join(Path(), "config.json")
}

func DefaultConfig() *Config {
	return &Config{
		Path:           Path(),
		DefaultVersion: "",
	}
}

func Read(path string) (*Config, error) {
	cfg := DefaultConfig()

	file, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		return cfg, nil
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		return nil, err
	}

	// if TFMGR_PATH is set, it should override the config
	if envPath, hasPath := os.LookupEnv(EnvPath); hasPath {
		cfg.Path = envPath
	}

	return cfg, nil
}

func Get() (*Config, error) {
	return Read(ConfigFile())
}

func (c *Config) Write() error {
	output, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	outputFile := ConfigFile()
	outputDir := filepath.Dir(outputFile)

	if _, err := os.Stat(outputDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		err = os.MkdirAll(outputDir, 0644)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(outputFile, output, 0644)
}
