package jsonml

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

const debug = false

type parseError error

type tokenError error

func IsParseError(err error) bool {
	_, answer := err.(parseError)
	return answer
}

func IsTokenError(err error) bool {
	_, answer := err.(tokenError)
	return answer
}

func PrintErrorTrace(err error) {
	fmt.Println(err.Error())
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s:%d\n", f, f)
		}
	}
}

func newParseError(expected string, actual json.Token) error {
	return parseError(errors.Errorf("Expected '%s' but got '%s'", expected, actual))
}

func newTokenError(err error) error {
	if err == io.EOF {
		return errors.New("Unexpected end of input stream")
	}

	if err != nil {
		return tokenError(errors.WithStack(err))
	}

	return nil
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
