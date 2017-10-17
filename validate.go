package main

import (
	"errors"
)

func ValidateArgs(args []string) error {
	argsLength := len(args)
	if argsLength != 1 {
		return errors.New("missing argument")
	}

	return nil
}
