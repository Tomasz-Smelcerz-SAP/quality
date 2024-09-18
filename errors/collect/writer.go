package main

import (
	"fmt"
	"strings"
)

// MyWriter is a simple WriteCloser that can be used for testing purposes.
type MyWriter struct {
	closeErr       error
	writeErr       error
	magicString    string
	respondWithErr bool
}

func NewWriter() *MyWriter {
	return &MyWriter{}
}

func (f *MyWriter) WithCloseError(err error) *MyWriter {
	f.closeErr = err
	return f
}

func (f *MyWriter) WithWriteError(err error) *MyWriter {
	f.writeErr = err
	return f
}

func (f *MyWriter) WithWriteErrorOnString(err error, s string) *MyWriter {
	f.writeErr = err
	f.magicString = s
	return f
}

func (f *MyWriter) Close() error {
	fmt.Println("[MyWriter]: Close()")

	if f.closeErr != nil {
		return f.closeErr
	}
	return nil
}

func (f *MyWriter) Write(p []byte) (n int, err error) {
	astr := string(p)

	fmt.Printf("[MyWriter]: Write(\"%s\")\n", astr)

	if len(f.magicString) > 0 && strings.Contains(astr, f.magicString) {
		f.respondWithErr = true
	}

	if f.respondWithErr {
		if f.writeErr != nil {
			return 0, f.writeErr
		}
	}

	return len(p), nil
}
