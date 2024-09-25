package main

import (
	"errors"
	"fmt"

	sgn "github.com/Tomasz-Smelcerz-SAP/quality/errors/signatures"
	dio "github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
)

var (
	ErrWriting = fmt.Errorf("error writing")
	ErrClosing = fmt.Errorf("error closing")
	ErrReading = fmt.Errorf("error reading")
)

func printErr(err error) {
	fmt.Println(err)
}

func main() {
	si := dio.NewStringReader([]string{"To", "be", "or", "not", "to", "be", "that", "is", "the", "question"}).
		WithErrorOn(ErrReading, "that")

	wCount, err := idiomaticWrite(si)
	fmt.Println("Written", wCount, "characters")
	if err != nil {
		fmt.Println("Errors during write:")
		errList := sgn.AsList(err)
		for _, e := range errList {
			fmt.Println("-", e)
		}
	}
}

func idiomaticWrite(si dio.StringIterator) (res int, err error) {
	w := dio.NewWriter().
		WithCloseError(ErrClosing)

	defer func() {
		wErr := w.Close()
		if wErr != nil {
			err = errors.Join(err, wErr)
		}

	}()

	totalCharsWritten := 0
	for si.HasNext() {
		nextStr, rErr := si.Next()
		if rErr != nil {
			return totalCharsWritten, rErr
		}
		count, wErr := w.Write([]byte(nextStr))
		totalCharsWritten += count
		if wErr != nil {
			return totalCharsWritten, wErr
		}
	}
	return totalCharsWritten, nil
}
