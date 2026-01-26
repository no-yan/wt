package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateName validates workspace/branch names for both Git and jj.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if strings.ContainsAny(name, ";&|`$(){}[]<>") {
		return fmt.Errorf("invalid characters in name: %q", name)
	}

	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("name cannot start with dash: %q", name)
	}

	return nil
}

// ValidatePath validates repository paths.
func ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute: %q", path)
	}

	return nil
}
