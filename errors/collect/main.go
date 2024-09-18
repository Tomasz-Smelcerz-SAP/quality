package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/Tomasz-Smelcerz-SAP/errs"
)

var ErrWriting = fmt.Errorf("error writing")
var ErrReading = fmt.Errorf("error reading")

func main() {
	fmt.Println("Let's go!")
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func run() error {

	reader := &stringReader{str: []string{"To", "be", "or", "not", "to", "be", "that", "is", "the", "question"}}

	clctr := errs.SimpleClctr{}
	totalCharsWritten, _ := safeWrite(reader, singletonWriter(), &clctr)
	fmt.Printf("Total characters written: %d\n", totalCharsWritten)
	if clctr.HasErrors() {
		return clctr.Errors()
	}

	return nil
}

// safeWrite reads from a stringIterator and writes to a writer. It returns the total number of characters written.
// It is using errs.Collector for error handling.
func safeWrite(reader stringIterator, gwf getWriterFn, c errs.Collector) (int, errs.Ignore) {
	// Acquire a writer and ensure it is closed after the function returns.
	writer := gwf()
	defer c.CollectF(writer.Close) // one-liner! No named return value needed.

	var totalWrittenCnt int

	for reader.hasNext() {
		strVal, rErr := reader.next()
		if rErr != nil {
			return totalWrittenCnt, c.Collect(fmt.Errorf("%w: %w", ErrReading, rErr))
		}
		writtenCount, wErr := writer.Write([]byte(strVal))
		totalWrittenCnt += writtenCount
		if wErr != nil {
			return totalWrittenCnt, c.Collect(fmt.Errorf("%w: %w", ErrWriting, wErr))
		}
	}

	return totalWrittenCnt, c.Collect(nil)
}

type getWriterFn func() io.WriteCloser

func singletonWriter() getWriterFn {
	created := false
	return func() io.WriteCloser {
		if created {
			panic("Writer already created")
		}
		res := NewWriter().WithCloseError(errors.New("cannot close the writer: I am too lazy"))
		created = true
		return res
	}
}

// stringIterator is an interface for reading strings.
type stringIterator interface {
	next() (string, error)
	hasNext() bool
}

type stringReader struct {
	index int
	str   []string
}

func (r *stringReader) next() (string, error) {
	if r.index >= len(r.str) {
		return "", io.EOF
	}
	s := r.str[r.index]
	r.index++
	return s, nil
}

func (r *stringReader) hasNext() bool {
	return r.index < len(r.str)
}
