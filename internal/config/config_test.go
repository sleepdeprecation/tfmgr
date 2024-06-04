package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sleepdeprecation/tfmgr/internal/config"
)

func TestConfig(t *testing.T) {
	t.Run("Test Path()", func(t *testing.T) {
		t.Setenv("HOME", "/some/home/dir")

		t.Run("With HOME set", func(t *testing.T) {
			expected := "/some/home/dir/.tfmgr"
			actual := config.Path()
			if actual != expected {
				t.Errorf("Config path is `%s`, should be `%s`", actual, expected)
			}
		})

		t.Run("With TFMGR_PATH set", func(t *testing.T) {
			t.Setenv("TFMGR_PATH", "/opt/tfmgr")

			expected := "/opt/tfmgr"
			actual := config.Path()
			if actual != expected {
				t.Errorf("Config path is `%s`, should be `%s`", actual, expected)
			}
		})
	})

	t.Run("Test ConfigFile()", func(t *testing.T) {
		t.Setenv("HOME", "/some/home/dir")

		t.Run("With HOME set", func(t *testing.T) {
			expected := "/some/home/dir/.tfmgr/config.json"
			actual := config.ConfigFile()
			if actual != expected {
				t.Errorf("Config file is `%s`, should be `%s`", actual, expected)
			}
		})

		t.Run("With TFMGR_CONFIG set", func(t *testing.T) {
			t.Setenv("TFMGR_CONFIG", "/opt/tfmgr/config.json")

			expected := "/opt/tfmgr/config.json"
			actual := config.ConfigFile()
			if actual != expected {
				t.Errorf("Config file is `%s`, should be `%s`", actual, expected)
			}
		})
	})

	t.Run("Test Get()", func(t *testing.T) {
		t.Run("With TFMGR_PATH set", func(t *testing.T) {
			configDir, err := os.MkdirTemp("", "tfmgr")
			if err != nil {
				panic(err)
			}

			t.Cleanup(func() {
				os.RemoveAll(configDir)
			})

			t.Setenv("TFMGR_PATH", configDir)

			t.Run("Without a config file", func(t *testing.T) {
				cfg, err := config.Get()
				if err != nil {
					t.Errorf("Error getting config file: %s", err.Error())
				}

				actualPath := cfg.Path
				expectedPath := configDir
				if actualPath != expectedPath {
					t.Errorf("Config value `path` is `%s`, should be `%s`", actualPath, expectedPath)
				}
			})

			t.Run("With a config file", func(t *testing.T) {
				err := os.WriteFile(
					filepath.Join(configDir, "config.json"),
					[]byte(`{"path": "/some/other/path"}`),
					0644)
				if err != nil {
					t.Errorf("Error writing config file: %s", err.Error())
				}

				cfg, err := config.Get()
				actualPath := cfg.Path
				expectedPath := configDir // because TFMGR_PATH is set, it overrides the config file
				if actualPath != expectedPath {
					t.Errorf("Config value `path` is `%s`, should be `%s`", actualPath, expectedPath)
				}
			})
		})
	})

	t.Run("Test Config.Write()", func(t *testing.T) {
		configDir, err := os.MkdirTemp("", "tfmgr")
		if err != nil {
			panic(err)
		}

		t.Cleanup(func() {
			os.RemoveAll(configDir)
		})

		t.Setenv("TFMGR_PATH", configDir)

		cfg := config.DefaultConfig()
		cfg.DefaultVersion = "1.0.0"

		err = cfg.Write()
		if err != nil {
			t.Errorf("Error writing config file: %s", err.Error())
		}

		read, err := config.Get()
		if err != nil {
			t.Errorf("Error reading config: %s", err.Error())
		}
		if read.Path != configDir {
			t.Errorf("Written configuration path is %s, should be %s", read.Path, configDir)
		}
		if read.DefaultVersion != "1.0.0" {
			t.Errorf("Written configuration default version is %s, should be %s", read.DefaultVersion, "1.0.0")
		}
	})
}
