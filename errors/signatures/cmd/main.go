package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/samber/mo"

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
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <number> <number>")
		return
	}

	firstNum := properAtoi(os.Args[1])
	secondNum := properAtoi(os.Args[2])

	if firstNum.IsRight() {
		printErr(firstNum.MustRight())
		return
	}
	if secondNum.IsRight() {
		printErr(secondNum.MustRight())
		return
	}

	if firstNum.IsLeft() && secondNum.IsLeft() {
		fmt.Printf("%s + %s = %d\n", os.Args[1], os.Args[2], firstNum.MustLeft()+secondNum.MustLeft())
	}

	si := dio.NewStringReader([]string{"To", "be", "or", "not", "to", "be", "that", "is", "the", "question"}).WithErrorOn(ErrReading, "that")

	res := safeWrite(si)
	fmt.Println("Written", res.Get(), "bytes")
	if res.HasError() {
		fmt.Println("Error writing to file:", res.GetError())
	}

	e1 := errors.New("error 1")
	e2 := errors.New("error 2")
	e3 := errors.New("error 3")
	e4 := errors.New("error 4")
	e5 := errors.New("error 5")

	errs := []error{e2, e3, e4, e5}

	resErr := e1
	for _, e := range errs {
		resErr = errors.Join(resErr, e)
	}

	fmt.Println("================================================================================")
	flatErrList := sgn.Flatten(resErr)
	for i, e := range flatErrList {
		fmt.Println(i, " -> ", e)
	}
	fmt.Println("================================================================================")
	resErrList := sgn.AsList(resErr)
	for i, e := range resErrList {
		fmt.Println(i, " -> ", e)
	}
}

// properAtoi is a wrapper around strconv.Atoi that returns either an int or an error
func properAtoi(s string) mo.Either[int, sgn.MandatoryError] {
	i, err := strconv.Atoi(s)
	if err != nil {
		return mo.Right[int, sgn.MandatoryError](sgn.NewMandatoryError(err))
	}
	return mo.Left[int, sgn.MandatoryError](i)
}

func safeWrite(si dio.StringIterator) (res sgn.IntWithOptErr) {
	w := NewExplicitWriter(dio.NewWriter())
	defer func() {
		res = res.Join(w.Close())
	}()

	totalCharsWritten := 0
	for si.HasNext() {
		nextStr, err := si.Next()
		if err != nil {
			return sgn.NewIntWithOptErr(totalCharsWritten, err)
		}
		res := w.Write([]byte(nextStr))
		totalCharsWritten += res.Get()
		if res.HasError() {
			return sgn.NewIntWithOptErr(totalCharsWritten, err)
		}
	}
	return sgn.NewIntWithOptErr(totalCharsWritten, nil)
}

func NewExplicitWriter(wc io.WriteCloser) explicitWriter {
	return explicitWriter{wc: wc}
}

// explicitWriter is a wrapper around io.WriteCloser that diverges from the standard Go error return convention.
// Instead of returning a value and an error, it returns a single value of type IntWithOptErr.
type explicitWriter struct {
	wc io.WriteCloser
}

func (tw explicitWriter) Write(p []byte) sgn.IntWithOptErr {
	return sgn.NewIntWithOptErr(tw.wc.Write(p))
}

func (tw explicitWriter) Close() error {
	return tw.wc.Close()
}
