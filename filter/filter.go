package filter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
)

func MatchAny(path string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		match, err := doublestar.Match(pattern, filepath.ToSlash(path))
		if err != nil {
			return false, fmt.Errorf("Error matching pattern '%s': %s", pattern, err)
		} else if match {
			return true, nil
		}
	}
	return false, nil
}

// File determines if a file should be included. Returns a cleaned path relative
// to the root, or the empty string if the file should be skipped.
func File(
	path string,
	includePatterns []string,
	excludePatterns []string,
) (string, error) {
	cleanpath := path
	if excluded, err := MatchAny(cleanpath, excludePatterns); err != nil {
		return "", err
	} else if excluded {
		return "", nil
	}
	if included, err := MatchAny(cleanpath, includePatterns); err != nil {
		return "", err
	} else if included {
		return cleanpath, nil
	}
	return "", nil
}

// Files filters an array of files using filter.File.
func Files(
	files []string,
	includePatterns []string,
	excludePatterns []string,
) ([]string, error) {
	ret := []string{}
	for _, file := range files {
		path, err := File(file, includePatterns, excludePatterns)
		if err != nil {
			continue
		}
		if path != "" {
			ret = append(ret, path)
		}
	}
	return ret, nil
}

// SplitPattern splits a pattern into a root directory and a trailing pattern
// specifier.
func SplitPattern(pattern string) (string, string) {
	base := pattern
	trail := ""

	split := strings.IndexAny(pattern, "*{}?[]")
	if split >= 0 {
		base = pattern[:split]
		trail = pattern[split:]
	}
	return base, trail
}
