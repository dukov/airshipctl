package document

import (
	"fmt"
)

// ErrDocNotFound returned if desired document not found
type ErrDocNotFound struct {
	Annotation string
	Kind       string
}

func (e ErrDocNotFound) Error() string {
	return fmt.Sprintf("Document annotated by %s with Kind %s not found", e.Annotation, e.Kind)
}

// ErrWrongRenderArgs returned if raw filter option is specified along
// with annotation or kind or label or apiversion
type ErrWrongRenderArgs struct{}

func (e ErrWrongRenderArgs) Error() string {
	return "Can not use raw filter argument with annotation or kind or label or apiversion"
}

// ErrBadRenderArgFormat returned if rendering keys have worng format
type ErrBadRenderArgFormat struct {
	Arg string
}

func (e ErrBadRenderArgFormat) Error() string {
	return fmt.Sprintf("Wrong format for %s. Expected format is 'someKey=someValue'", e.Arg)
}
