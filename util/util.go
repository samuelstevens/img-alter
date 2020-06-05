package util

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// HashFile returns a SHA1 hash of a filepath (not cryptographically secure)
func HashFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()

	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}

func GreatestCommonAbs(general string, specific string) (string, error) {
	if !filepath.IsAbs(general) {
		return "", fmt.Errorf("general: %s must be an absolute path", general)
	}

	if !filepath.IsAbs(specific) {
		return "", fmt.Errorf("specific: %s must be an absolute path", specific)
	}

	generalParts := splitPath(general)
	specificParts := splitPath(specific)

	bestOption := []string{}

	for i := 0; i < len(generalParts); i++ {
		if generalParts[i] == specificParts[i] {
			bestOption = generalParts[:i+1]
		}
	}

	return filepath.FromSlash("/" + strings.Join(bestOption, "/")), nil
}

func root(path string) string {
	if path[0] == '.' {
		return "."
	}

	if filepath.Dir(path) == "." {
		return path
	}

	if filepath.Dir(path) == "/" {
		return path
	}

	return root(filepath.Dir(path))
}

func splitPath(path string) []string {
	split := strings.Split(filepath.ToSlash(path), "/")

	result := []string{}

	for _, s := range split {
		if s != "" {
			result = append(result, s)
		}
	}
	// strip empty strings
	return result
}

// Dir returns the path's directory, checking if it's already a directory. The path must already exist.
func Dir(path string) (string, error) {
	info, err := os.Stat(path)

	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return path, nil
	}

	return filepath.Dir(path), nil
}

// Join two lists of directories, assuming the first is an absolute path
func Join(abs []string, other []string) string {
	absPath := "/" + strings.Join(abs, "/")
	otherPath := strings.Join(other, "/")

	totalPath := absPath + "/" + otherPath

	return filepath.Clean(filepath.FromSlash(totalPath))
}

// MakeAbsRelativeTo takes an absolute path and an arbitrary path and finds the arbitrary path's absolute equivalent by looking for it within the absolute path.
// Depends on the otherPath being a real file on disk.
func MakeAbsRelativeTo(absPath string, otherPath string) (string, error) {
	absPath = filepath.Clean(absPath)
	otherPath = filepath.Clean(otherPath)

	absDir, err := Dir(absPath)

	if err != nil {
		return "", err
	}

	if !filepath.IsAbs(absDir) {
		return "", fmt.Errorf("%s must be absolute", absDir)
	}

	absParts := splitPath(absDir)

	otherParts := splitPath(otherPath)

	for i := len(absParts); i >= 0; i-- {
		testPath := Join(absParts[:i], otherParts)

		_, err := os.Stat(testPath)

		if err != nil {
			continue
		}

		return testPath, nil
	}

	return "", fmt.Errorf("%s not found on disk", otherPath)
}

func ExpandUserDirectory(path string) string {
	homedir, err := os.UserHomeDir()

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"

		if err != nil {
			return path
		}

		return homedir

	} else if strings.HasPrefix(path, "~/") {
		if err != nil {
			return path
		}

		return filepath.Join(homedir, path[2:])
	}

	return path
}

// StringSet is a set of strings
type StringSet map[string]struct{}

// Exists sentinel value
var Exists = struct{}{}

// NewStringSet makes a new StringSet from a slice of strings
func NewStringSet(elements []string) *StringSet {
	set := StringSet{}
	for _, elem := range elements {
		set[elem] = Exists
	}

	return &set
}

// Contains checks if a string is in a StringSet
func (set *StringSet) Contains(element string) bool {
	_, ok := (*set)[element]

	return ok
}

// Empty checks if a StringSet is empty
func (set *StringSet) Empty() bool {
	return len(*set) == 0
}

// Remove an element from a set
func (set *StringSet) Remove(element string) {
	delete(*set, element)
}

type KeyError struct {
	key string
}

func (e *KeyError) Error() string {
	return fmt.Sprintf("key error: %s not in set", e.key)
}
