package main

import (
	"testing"
)

var args []string

func TestValidateArgsSuccess(t *testing.T) {
	args = []string {
		"hoge",
	}

	err := ValidateArgs(args)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
}

func TestValidateArgsError(t *testing.T) {

}
