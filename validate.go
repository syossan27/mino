package main

import (
	"errors"
	"os"
)

func ValidateArgs() error {
	argsLength := len(os.Args)
	if argsLength != 2 {
		return errors.New("missing argument")
	}

	return nil
}
