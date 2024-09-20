package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	tuple "github.com/barweiss/go-tuple"
	"github.com/samber/mo"

	dio "github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
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

}

// properAtoi is a wrapper around strconv.Atoi that returns either an int or an error
func properAtoi(s string) mo.Either[int, MandatoryError] {
	i, err := strconv.Atoi(s)
	if err != nil {
		return mo.Right[int, MandatoryError](NewMandatoryError(err))
	}
	return mo.Left[int, MandatoryError](i)
}

func safeWrite(si dio.StringIterator) (res EitherIntOrErr) {
	w := NewTypedWriter(dio.NewWriter())
	defer func() {
		if err := w.Close(); err != nil {
			if res.IsLeft() {
				res = EitherIntOrErr{mo.Right[int, ErrWithInt](
					ErrWithInt(tuple.T2[MandatoryError, int]{
						V1: NewMandatoryError(err),
						V2: res.MustLeft(),
					}),
				)}
			} else {
				res = EitherIntOrErr{mo.Right[int, ErrWithInt](
					ErrWithInt(tuple.T2[MandatoryError, int]{
						V1: NewMandatoryError(err),
						V2: res.MustRight().V2,
					}),
				)}
			}
		}
	}()

	reserr := w.Write([]byte("Hello, World!"))

	if reserr.IsLeft() {
		fmt.Printf("Written %d bytes\n", reserr.MustLeft())
		return
	}

	fmt.Println("Error writing to file:", reserr.MustRight().V1.Error())
	fmt.Println("number of bytes returned: ", reserr.MustRight().V2)

	return
}

func NewTypedWriter(wc io.WriteCloser) typedWriter {
	return typedWriter{wc: wc}
}

type typedWriter struct {
	wc io.WriteCloser
}

func (tw typedWriter) Write(p []byte) EitherIntOrErr {
	bCount, err := tw.wc.Write(p)
	if err != nil {
		return EitherIntOrErr{mo.Right[int, ErrWithInt](
			ErrWithInt(tuple.T2[MandatoryError, int]{
				V1: NewMandatoryError(err),
				V2: bCount,
			}),
		)}
	}
	return EitherIntOrErr{mo.Left[int, ErrWithInt](bCount)}
}

func (tw typedWriter) Close() error {
	return tw.wc.Close()
}
