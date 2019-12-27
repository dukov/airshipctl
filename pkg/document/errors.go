package document

import (
	"fmt"
)

// ErrDocNotFound returned if desired document not found
type ErrDocNotFound struct {
	Label string
	Kind  string
}

func (e ErrDocNotFound) Error() string {
	return fmt.Sprintf("Document labeled by %s with Kind %s not found", e.Label, e.Kind)
}
