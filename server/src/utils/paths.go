package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

type ProjectPaths struct {
	WorkDir     string // Current working directory
	ExecDir     string // Directory containing the executable
	ProjectRoot string // Root directory of the project
}

// GetProjectPaths returns various useful project paths
func GetProjectPaths() (*ProjectPaths, error) {
	paths := &ProjectPaths{}
	var err error

	// Get working directory
	paths.WorkDir, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	// Get executable directory
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	paths.ExecDir = filepath.Dir(execPath)

	// Get project root (using go.mod as marker)
	_, b, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(b)

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			paths.ProjectRoot = currentDir
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// If we can't find go.mod, use working directory as fallback
			paths.ProjectRoot = paths.WorkDir
			break
		}
		currentDir = parent
	}

	return paths, nil
}
