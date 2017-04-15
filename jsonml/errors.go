package jsonml

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

const debug = false

func IsParseError(err error) bool {
	if de, ok := err.(decodeError); ok && de.parse {
		return true
	}

	return false
}

func PrintErrorTrace(err error) {
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s:%d\n", f, f)
		}
	}
}

func parseError(expected string, actual json.Token) error {
	return decodeError{parse: true, err: errors.Errorf("Expected %s, but got: %s", expected, actual)}
}

func tokenError(err error) error {
	if err == io.EOF {
		return errors.New("Unexpected end of input stream")
	}

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
