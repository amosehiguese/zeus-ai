package cobrautil

import (
	"fmt"
	"os"
)

const (
	ExitSuccess = iota
	ExitError
)

func ExitWithError(code int, err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(code)
}
