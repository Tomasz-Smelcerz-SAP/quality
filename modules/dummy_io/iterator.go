package dummy_io

import (
	"io"
)

// StringIterator is an interface for reading strings.
type StringIterator interface {
	Next() (string, error)
	HasNext() bool
}

type StringReader struct {
	index int
	str   []string
}

func NewStringReader(str []string) *StringReader {
	return &StringReader{str: str}
}

func (r *StringReader) Next() (string, error) {
	if r.index >= len(r.str) {
		return "", io.EOF
	}
	s := r.str[r.index]
	r.index++
	return s, nil
}

func (r *StringReader) HasNext() bool {
	return r.index < len(r.str)
}
