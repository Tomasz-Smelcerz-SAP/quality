package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Tomasz-Smelcerz-SAP/errs"
	"github.com/Tomasz-Smelcerz-SAP/quality/modules/testwriter"
)

var ErrWriting = fmt.Errorf("error writing")
var ErrReading = fmt.Errorf("error reading")

func main() {
	fmt.Println("Let's write something!")

	runFunc := runClassic
	if len(os.Args) > 1 && os.Args[1] == "collect" {
		runFunc = runCollect
	}

	reader := &stringReader{str: []string{"To", "be", "or", "not", "to", "be", "that", "is", "the", "question"}}

	totalCharsWritten, err := runFunc(reader)
	fmt.Printf("Total characters written: %d\n", totalCharsWritten)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// writeClassic reads from a stringIterator and writes to a writer. It returns the total number of characters written.
// It is using classic golang error handling style.
func writeClassic(reader stringIterator, gwf getWriterFn) (totalWrittenCnt int, err error) {
	// Acquire a writer. This is equivalent to, say, opening a File.
	writer := gwf()
	// Ensure the writer is closed.
	defer func() {
		if cErr := writer.Close(); cErr != nil {
			err = errors.Join(cErr)
		}
	}()

	for reader.hasNext() {
		strVal, rErr := reader.next()
		if rErr != nil {
			err = fmt.Errorf("%w: %w", ErrReading, rErr) // definition in the separate line for better readability
			return
		}
		writtenCount, wErr := writer.Write([]byte(strVal))
		totalWrittenCnt += writtenCount
		if wErr != nil {
			err = fmt.Errorf("%w: %w", ErrWriting, wErr)
			return
		}
	}
	return
}

// writeCollect reads from a stringIterator and writes to a writer. It returns the total number of characters written.
// It is using errs.Collector for error handling.
func writeCollect(reader stringIterator, gwf getWriterFn, c errs.Collector) (int, errs.Ignore) {
	// Acquire a writer. This is equivalent to, say, opening a File.
	writer := gwf()
	// Ensure the writer is closed. Note it's a one-liner and no named return values are used.
	defer c.CollectF(writer.Close)

	var totalWrittenCnt int

	for reader.hasNext() {
		strVal, rErr := reader.next()
		if rErr != nil {
			err := fmt.Errorf("%w: %w", ErrReading, rErr) // definition in the separate line for better readability
			return totalWrittenCnt, c.Collect(err)
		}
		writtenCount, wErr := writer.Write([]byte(strVal))
		totalWrittenCnt += writtenCount
		if wErr != nil {
			err := fmt.Errorf("%w: %w", ErrWriting, wErr)
			return totalWrittenCnt, c.Collect(err)
		}
	}

	// the second value is equivalent to returning nil error in the idiomatic Go code
	return totalWrittenCnt, c.Collect(nil)
}

func runClassic(sIter stringIterator) (int, error) {
	fmt.Println("Classic approach is the best!")

	totalCharsWritten, err := writeClassic(sIter, singletonWriter())
	return totalCharsWritten, err
}

func runCollect(sIter stringIterator) (int, error) {
	fmt.Println("Classic approach isn't the best!")

	clctr := errs.SimpleClctr{}
	totalCharsWritten, _ := writeCollect(sIter, singletonWriter(), &clctr)
	return totalCharsWritten, clctr.Errors()
}

////////////////////// Helper stuff

type getWriterFn func() io.WriteCloser

func singletonWriter() getWriterFn {
	created := false
	return func() io.WriteCloser {
		if created {
			panic("Writer already created")
		}
		res := testwriter.NewWriter().
			WithCloseError(errors.New("cannot close the writer: I am too lazy")).
			WithWriteErrorOnString(ErrWriting, "question")
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
