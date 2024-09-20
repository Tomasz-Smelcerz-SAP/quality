package main

import (
	"errors"
	"fmt"

	dio "github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
)

func runErrorList(sIter dio.StringIterator) (int, error) {
	fmt.Println("Having a list of errors is the best way!\n")

	totalWrittenCnt, errs := writeErrorList(sIter, SingletonWriter())
	return totalWrittenCnt, errors.Join(errs...)
}

// writeErrorList reads from a StringIterator and writes to a Writer. It returns the total number of characters written, and an optional list of errors.
func writeErrorList(reader dio.StringIterator, gwf GetWriterFn) (totalWrittenCnt int, errs []error) {
	// Acquire a writer. This is equivalent to, say, opening a File.
	writer := gwf()
	// Ensure the writer is closed.
	defer func() {
		cErr := writer.Close()
		if cErr != nil {
			errs = append(errs, cErr)
		}
	}()

	for reader.HasNext() {
		strVal, rErr := reader.Next()
		if rErr != nil {
			err := fmt.Errorf("%w: %w", ErrReading, rErr)
			errs = append(errs, err)
			return
		}
		writtenCount, wErr := writer.Write([]byte(strVal))
		totalWrittenCnt += writtenCount
		if wErr != nil {
			err := fmt.Errorf("%w: %w", ErrWriting, wErr)
			errs = append(errs, err)
			return
		}
	}
	return
}
