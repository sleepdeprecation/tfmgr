package manager

import (
	"os"
	"path/filepath"
	"strings"
)

type VersionDetector struct {
	path string
}

func NewDetector(path string) *VersionDetector {
	return &VersionDetector{
		path: path,
	}
}

func (d *VersionDetector) HasVersionFile() (bool, error) {
	stat, err := os.Stat(filepath.Join(d.path, ".terraform-version"))
	if err != nil {
		return false, nil
	}

	return true, nil
}

// CheckVersionFile checks if a .terraform-version file exists in the current tree, and if it does, it returns the version specified in the file.
func (d *VersionDetector) CheckVersionFile() (string, error) {
	path := filepath.Join(d.path, ".terraform-version")
	body, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(body))
	return version, nil
}
