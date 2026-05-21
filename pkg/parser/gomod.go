package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadModulePath walks up from dir until it finds a go.mod file, then returns
// the module path declared in it. Returns an error if no go.mod is found.
func ReadModulePath(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(dir, "go.mod")
		f, err := os.Open(candidate)
		if err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, "module ") {
					return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
				}
			}
			return "", fmt.Errorf("no module declaration found in %s", candidate)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root
		}
		dir = parent
	}

	return "", fmt.Errorf("no go.mod found in %s or any parent directory", dir)
}
