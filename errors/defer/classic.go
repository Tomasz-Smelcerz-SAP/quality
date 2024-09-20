package main

import (
	"errors"
	"fmt"

	dio "github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
)

func runClassic(sIter dio.StringIterator) (int, error) {
	fmt.Println("Classic approach is the best!\n")

	return writeClassic(sIter, SingletonWriter())
}

// writeClassic reads from a StringIterator and writes to a Writer. It returns the total number of characters written.
// It is using classic golang error handling style.
func writeClassic(reader dio.StringIterator, gwf GetWriterFn) (totalWrittenCnt int, err error) {
	// Acquire a writer. This is equivalent to, say, opening a File.
	writer := gwf()
	// Ensure the writer is closed.
	defer func() {
		cErr := writer.Close()
		err = errors.Join(err, cErr) //NOTE: Initially I made a mistake here: `err = errors.Join(cErr)` (without the second argument)
	}()

	for reader.HasNext() {
		strVal, rErr := reader.Next()
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
