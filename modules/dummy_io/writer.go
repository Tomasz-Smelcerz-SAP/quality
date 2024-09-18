package dummy_io

import (
	"fmt"
	"strings"
)

// MyWriter is a simple WriteCloser that can be used for testing purposes.
// It can be configured to return an error on Close() and/or Write().
type MyWriter struct {
	closeErr       error  // error to return on Close()
	writeErr       error  // error to return on Write()
	magicString    string // if this string is written, respond with a configured writeErr
	respondWithErr bool
}

func NewWriter() *MyWriter {
	return &MyWriter{}
}

func (f *MyWriter) WithCloseError(err error) *MyWriter {
	f.closeErr = err
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
	if len(p) == 0 {
		return 0, nil
	}

	aStr := string(p)

	fmt.Printf("[MyWriter]: Write(\"%s\")\n", aStr)

	if len(f.magicString) > 0 && strings.Contains(aStr, f.magicString) {
		f.respondWithErr = true
	}

	if f.respondWithErr {
		if f.writeErr != nil {
			return len(p) - 1, f.writeErr // pretend we wrote all but the last byte
		}
	}

	return len(p), nil
}
