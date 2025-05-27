package analyzer

import "fmt"

type NonExistingFileError struct {
	Path string
	Err  error
}

func (e *NonExistingFileError) Error() string {
	return fmt.Sprintf("non-existing file %s: %v", e.Path, e.Err)
}

func (e *NonExistingFileError) Unwrap() error {
	return e.Err
}
