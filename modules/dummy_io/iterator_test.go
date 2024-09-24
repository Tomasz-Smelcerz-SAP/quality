package dummy_io_test

import (
	"errors"
	"io"
	"testing"

	"github.com/Tomasz-Smelcerz-SAP/quality/modules/dummy_io"
	"github.com/stretchr/testify/assert"
)

func TestStaticReader(t *testing.T) {
	// Given
	str := []string{"a", "b", "c"}
	reader := dummy_io.NewStringReader(str)

	// When
	var res []string
	for reader.HasNext() {
		val, err := reader.Next()
		assert.NoError(t, err)
		res = append(res, val)
	}

	// Then
	assert.Equal(t, str, res)
}

func TestStaticReaderWithError(t *testing.T) {
	// Given
	str := []string{"a", "b", "c"}
	fakeErr := errors.New("fake error")
	reader := dummy_io.NewStringReader(str).WithErrorOn(fakeErr, "c")

	// When
	var elems []string
	var actualErr error
	for reader.HasNext() {
		val, err := reader.Next()
		if err != nil && !errors.Is(err, io.EOF) {
			actualErr = err
			break
		}
		elems = append(elems, val)
	}

	// Then
	assert.Equal(t, []string{"a", "b"}, elems)
	assert.Equal(t, fakeErr, actualErr)
}
