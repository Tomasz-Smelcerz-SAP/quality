package main

import (
	"fmt"

	"github.com/Tomasz-Smelcerz-SAP/errs"
	dio "github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
)

func runCollect(sIter dio.StringIterator) (int, error) {
	fmt.Println("Collecting errors is better than returning them!\n")

	clctr := errs.SimpleClctr{}
	totalWrittenCnt, _ := writeCollect(sIter, SingletonWriter(), &clctr)
	return totalWrittenCnt, clctr.Errors()
}

// writeCollect reads from a StringIterator and writes to a Writer. It returns the total number of characters written.
// It is using errs.Collector for error handling.
func writeCollect(reader dio.StringIterator, gwf GetWriterFn, c errs.Collector) (int, errs.Ignore) {
	// Acquire a writer. This is equivalent to, say, opening a File.
	writer := gwf()
	// Ensure the writer is closed. Note it's a one-liner and no named return values are used.
	defer c.CollectF(writer.Close)

	var totalWrittenCnt int

	for reader.HasNext() {
		strVal, rErr := reader.Next()
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
