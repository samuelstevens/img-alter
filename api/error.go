package api

import "fmt"

// ConfidenceError indicates that Azure could not find a caption with a high confidence.
type ConfidenceError struct {
	Confidence float64
}

func (e *ConfidenceError) Error() string {
	return fmt.Sprintf("low confidence: %f", e.Confidence)
}

// NoLabelError indicates that Azure could not find any captions for the path.
type NoLabelError struct {
	Path string
}

func (e *NoLabelError) Error() string {
	return fmt.Sprintf("no descriptions found for %s", e.Path)
}
