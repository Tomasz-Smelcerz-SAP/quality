package main

import "errors"

func NewMandatoryError(e error) MandatoryError {
	return MandatoryError{e: e, msg: e.Error()}
}

// mandatoryError is a custom error type that is never nil
type MandatoryError struct {
	e   error
	msg string
}

func (e MandatoryError) Get() error {
	return e.e
}

func (e MandatoryError) Error() string {
	return e.msg
}

// Add returns a new instance of MandatoryError that joins the current error with another error.
func (me MandatoryError) Add(another error) MandatoryError {
	if another == nil {
		return me
	}

	var errList []error

	// Current error may either be a single error, or it may be a joined list of errors.
	// If it is a joined list of errors, we need to extract the individual errors, because errors.Join does not "flatten" the internal error list.
	unwrapper, ok := me.e.(unwrapper)
	if ok {
		errList = unwrapper.Unwrap()
	} else {
		errList = append(errList, me.e)
	}

	return NewMandatoryError(errors.Join(errList...))
}

type unwrapper interface {
	Unwrap() []error
}
