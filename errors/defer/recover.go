package main

import (
	"errors"
	"fmt"

	dio "github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
)

func runPanicRecover(sIter dio.StringIterator) (totalWrittenCnt int, err error) {

	pw := &panickingWriter{}

	defer func() {
		pVal := recover()

		var pErr error
		if pVal != nil {
			pErr = pVal.(error) //re-panic if not an error
		}

		totalWrittenCnt = pw.totalWrittenCnt
		err = errors.Join(pw.err, pErr)
	}()

	fmt.Println("Advanced error handling techniques for the win!\n")
	pw.writePanicRecover(sIter, SingletonWriter())
	return pw.totalWrittenCnt, err
}

type panickingWriter struct {
	totalWrittenCnt int
	err             error
}

// writePanicRecover reads from a StringIterator and writes to a Writer.
// It is using panic and recover to handle error during deferred Close() call.
func (pw *panickingWriter) writePanicRecover(reader dio.StringIterator, gwf GetWriterFn) {
	// Acquire a writer. This is equivalent to, say, opening a File.
	writer := gwf()
	// Ensure the writer is closed.
	defer func() {
		if cErr := writer.Close(); cErr != nil {
			panic(cErr)
		}
	}()

	for reader.HasNext() {
		strVal, rErr := reader.Next()
		if rErr != nil {
			pw.err = fmt.Errorf("%w: %w", ErrReading, rErr) // defined in the separate line for better readability
			return
		}
		writtenCount, wErr := writer.Write([]byte(strVal))
		pw.totalWrittenCnt += writtenCount
		if wErr != nil {
			pw.err = fmt.Errorf("%w: %w", ErrWriting, wErr) // defined in the separate line for better readability
			return
		}
	}
	return
}
