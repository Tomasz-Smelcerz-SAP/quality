package main

import (
	"errors"
	"fmt"

	dio "github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
)

func runDualError(sIter dio.StringIterator) (int, error) {
	fmt.Println("Being explicit about possible errors is the best solution!\n")

	totalWrittenCnt, err, cErr := writeDualError(sIter, SingletonWriter())
	return totalWrittenCnt, errors.Join(err, cErr)
}

// writeDualError reads from a StringIterator and writes to a Writer. It returns the total number of characters written.
// It is using a modified scheme of error handling: it returns two error values, a "standard" one and one dedicated for "Writer.Close()" error.
func writeDualError(reader dio.StringIterator, gwf GetWriterFn) (totalWrittenCnt int, err error, closeErr error) {
	// Acquire a writer. This is equivalent to, say, opening a File.
	writer := gwf()
	// Ensure the writer is closed.
	defer func() {
		closeErr = writer.Close()
	}()

	for reader.HasNext() {
		strVal, rErr := reader.Next()
		if rErr != nil {
			err = fmt.Errorf("%w: %w", ErrReading, rErr)
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
