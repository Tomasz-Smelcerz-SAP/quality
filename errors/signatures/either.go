package signatures

import "errors"

// IntWithOptErr is a custom error type representing an integer value with an optional error.
// It is useful to represent return values of functions like io.Writer.Write, which always return an integer value and may in addition return an error.
type IntWithOptErr struct {
	val int
	err error
}

func NewIntWithOptErr(val int, err error) IntWithOptErr {
	return IntWithOptErr{val: val, err: err}
}

func (i IntWithOptErr) HasError() bool {
	return i.err != nil
}

func (i IntWithOptErr) Get() int {
	return i.val
}

func (i IntWithOptErr) GetError() error {
	return i.err
}

// Error implements error.Error()
func (i IntWithOptErr) Error() string {
	return i.err.Error()
}

// Join returns a new instance of IntWithOptErr that joins the current error with another error.
func (i IntWithOptErr) Join(another error) IntWithOptErr {
	if another == nil {
		return i
	}

	if i.err == nil {
		return NewIntWithOptErr(i.val, another)
	}

	return NewIntWithOptErr(i.val, errors.Join(i.err, another))
}
