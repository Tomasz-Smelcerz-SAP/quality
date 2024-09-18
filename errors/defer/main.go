package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	dio "github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
	"github.com/Tomasz-Smelcerz-SAP/quality/modules/testwriter"
)

var ErrWriting = fmt.Errorf("error writing")
var ErrClosing = fmt.Errorf("error closing")
var ErrReading = fmt.Errorf("error reading")

func main() {
	fmt.Println("Let's write something!")

	// Decide which function to run based on the command line arguments.
	runFunc := runClassic
	if len(os.Args) > 1 && os.Args[1] == "collect" {
		runFunc = runCollect
	}

	// Given
	reader := dio.NewStringReader([]string{"To", "be", "or", "not", "to", "be", "that", "is", "the", "question"})

	// When
	totalCharsWritten, err := runFunc(reader)

	// Then
	fmt.Printf("Total characters written: %d\n", totalCharsWritten)
	if err != nil {
		fmt.Println("\nErrors encountered:")
		// Unwrap the error to get the original error.
		unwrapper, ok := err.(unwrapper)
		if ok {
			for _, e := range unwrapper.Unwrap() {
				fmt.Println("-->", e.Error())
			}
		} else {
			fmt.Println("-->", err.Error())
		}
	}
}

type unwrapper interface {
	Unwrap() []error
}

// Returns a function that returns a writer. Ensures the returned function panics if invoked more than once.
func SingletonWriter() dio.GetWriterFn {
	created := false
	return func() io.WriteCloser {
		if created {
			panic("Writer already created")
		}
		res := testwriter.NewWriter().
			WithCloseError(ErrClosing).
			WithWriteErrorOnString(errors.New("I am too lazy"), "question")
		created = true
		return res
	}
}
