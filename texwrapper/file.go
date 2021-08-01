package latex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// SourcePath gets a relative source directory below the provided root
// it evaluates symbolic links and prevents access outside this root
func (d *SourceDirectory) SourcePath(relpath string) (string, error) {
	naive := filepath.Join(d.sourcePath, relpath)

	followLinks, err := filepath.EvalSymlinks(naive)
	if err != nil {
		return "", fmt.Errorf("Failed to resolve %s: %w", relpath, err)
	}

	abspath, err := filepath.Abs(followLinks)
	if err != nil {
		return "", fmt.Errorf("Failed to get absolute %s: %w", followLinks, err)
	}

	if !strings.HasPrefix(abspath, d.sourcePath) {
		return "", fmt.Errorf("Access Error - %s is not below %s", relpath, d.sourcePath)
	}

	return naive, nil
}

// ListPath gets all the files relevent to the UI in the given directory
// (the subdirectories and *.tex files)
// It is expected that the provided directory has already gone through SourcePath
func (d *SourceDirectory) ListPath(directory string) ([]os.FileInfo, error) {
	info, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("Failed to read %s: %w", directory, err)
	}
	// filter out anything that isn't a .tex or directory
	removed := 0
	for i, file := range info {
		if !file.IsDir() && filepath.Ext(file.Name()) != ".tex" {
			removed++
		} else {
			info[i-removed] = file
		}
	}
	info = info[:len(info)-removed]
	return info, nil
}
