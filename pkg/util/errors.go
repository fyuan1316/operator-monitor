package util

import "strings"

type Errors struct {
	errs []error
}

func (e *Errors) Append(err error) {
	if _, ok := err.(error); ok {
		e.errs = append(e.errs, err)
	}
}
func (e *Errors) Empty() bool {
	return len(e.errs) == 0
}
func (e *Errors) Error() string {
	var errStrs []string
	for _, err := range e.errs {
		errStrs = append(errStrs, err.Error())
	}
	return strings.Join(errStrs, "\n")
}
func (e *Errors) ErrorOrNil() error {
	if e == nil {
		return nil
	}
	if len(e.errs) == 0 {
		return nil
	}
	return e
}
