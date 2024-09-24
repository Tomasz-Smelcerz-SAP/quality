package signatures

// IntWithOptErr is a custom error type representing an integer value with an optional error.
// It is useful to represent return values of functions like io.Writer.Write, which return an integer value along with an error.
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
