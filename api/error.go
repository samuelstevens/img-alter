package api

import (
	"errors"
	"fmt"
)

// ConfidenceError indicates that Azure could not find a caption with a high confidence.
type ConfidenceError struct {
	Confidence float64
}

func (e *ConfidenceError) Error() string {
	return fmt.Sprintf("low confidence: %f", e.Confidence)
}

// ErrorNoLabel indicates that Azure could not find any captions for the path.
var ErrorNoLabel = errors.New("no descriptions found")
