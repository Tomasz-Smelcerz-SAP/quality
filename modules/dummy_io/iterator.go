package dummy_io

import (
	"io"
)

// StringIterator is an interface for reading strings.
type StringIterator interface {
	Next() (string, error)
	HasNext() bool
}

type StaticReader struct {
	index   int
	str     []string
	errorOn string
	err     error
}

func NewStringReader(str []string) *StaticReader {
	return &StaticReader{str: str}
}

func (r *StaticReader) WithErrorOn(err error, val string) *StaticReader {
	r.err = err
	r.errorOn = val
	return r
}

func (r *StaticReader) Next() (string, error) {
	if r.index >= len(r.str) {
		return "", io.EOF
	}
	res := r.str[r.index]
	if r.err != nil && res == r.errorOn {
		return "", r.err
	}
	r.index++
	return res, nil
}

func (r *StaticReader) HasNext() bool {
	return r.index < len(r.str)
}
