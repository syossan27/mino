package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-colorable"
)

var (
	stdout = colorable.NewColorableStdout()
	stderr = colorable.NewColorableStderr()
)

func Success(msg string) {
	fmt.Fprintf(stdout, "\x1b[32m%s\x1b[0m\n", msg)
}

func Fatal(err error) {
	fmt.Fprintf(stderr, "\x1b[31mError: %s\x1b[0m\n", err)
	os.Exit(ExitCodeError)
}
