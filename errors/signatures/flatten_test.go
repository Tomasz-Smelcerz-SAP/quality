package signatures_test

import (
	"errors"
	"testing"

	sgn "github.com/Tomasz-Smelcerz-SAP/quality/errors/signatures"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	e1 := errors.New("error 1")
	e2 := errors.New("error 2")
	e3 := errors.New("error 3")
	e4 := errors.New("error 4")
	e5 := errors.New("error 5")

	errs := []error{e2, e3, e4, e5}

	//given
	resErr := e1
	for _, e := range errs {
		// joining in the loop results in nested joinErrors
		resErr = errors.Join(resErr, e)
	}
	errList := sgn.AsList(resErr)
	assert.Equal(t, 2, len(errList))
	assert.Equal(t, e5, errList[1]) // error 5 is the last in the list. All other errors are joined in the first element and cannot be easily accessed

	//when
	flatErrList := sgn.Flatten(resErr)
	assert.Equal(t, 5, len(flatErrList))
	assert.Equal(t, e1, flatErrList[0])
	assert.Equal(t, e2, flatErrList[1])
	assert.Equal(t, e3, flatErrList[2])
	assert.Equal(t, e4, flatErrList[3])
	assert.Equal(t, e5, flatErrList[4])
}
