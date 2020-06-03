package webpage

import "fmt"

// FileTypeError occurs when a WebPage doesn't get an .html file
type FileTypeError struct {
	path string
}

func (e *FileTypeError) Error() string {
	return fmt.Sprintf("%s is not an .html file", e.path)
}
