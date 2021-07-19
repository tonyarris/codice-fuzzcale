package main

import (
	"testing"
)

// test checkName
func TestValidName(t *testing.T) {
	// test valid name
	err := checkName("marco")
	if err != nil {
		t.Fatalf(`checkName("marco") results in error`)
	}
}

func TestInvalidName(t *testing.T) {
	err := checkName("123")
	if err == nil {
		t.Fatalf(`An invalid numerical name has been allowed`)
	}

	err = checkName("O'Hanlan")
	if err == nil {
		t.Fatalf(`An apostrophe has been allowed`)
	}
}
