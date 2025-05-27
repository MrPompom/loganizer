package checker

import "fmt"

type UnreachableURLError struct {
	URL string
	Err error
}

func (e *UnreachableURLError) Error() string {
	return fmt.Sprintf("unreachable URL %s: %v", e.URL, e.Err)
}

func (e *UnreachableURLError) Unwrap() error {
	return e.Err
}
