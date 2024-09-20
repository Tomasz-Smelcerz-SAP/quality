package main

import (
	tuple "github.com/barweiss/go-tuple"
	"github.com/samber/mo"
)

type ErrWithInt tuple.T2[MandatoryError, int]

// This is a struct because there is no working type aliases in Go...
type EitherIntOrErr struct {
	mo.Either[int, ErrWithInt]
}

// No function overloading in Go...
func NewEitherIntOrErrLeft(val int) EitherIntOrErr {
	return EitherIntOrErr{mo.Left[int, ErrWithInt](val)}
}

// No function overloading in Go...
func NewEitherIntOrErrRight(err error, i int) ErrWithInt {
	return ErrWithInt(tuple.T2[MandatoryError, int]{
		V1: NewMandatoryError(err),
		V2: i,
	})
}

func (e EitherIntOrErr) GetInt() int {
	if e.IsLeft() {
		return e.MustLeft()
	}
	return e.MustRight().V2
}

// AddErr returns a new instance of EitherIntOrErr that represents the given error and the current int value.
// If the existing instance already IsRight(), the new error is joined with the existing error.
func (ioe EitherIntOrErr) AddErr(err error) EitherIntOrErr {
	if err == nil {
		return ioe
	}

	intVal := ioe.GetInt()
	var newErr MandatoryError
	if ioe.IsLeft() {
		newErr = NewMandatoryError(err)
	} else {
		newErr = ioe.MustRight().V1.Add(err)
	}

	return EitherIntOrErr{mo.Right[int, ErrWithInt](ErrWithInt(tuple.T2[MandatoryError, int]{
		V1: newErr,
		V2: intVal,
	}))}
}
