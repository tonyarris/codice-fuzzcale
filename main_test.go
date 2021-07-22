package main

import (
	"testing"
)

// test checkName
func TestValidName(t *testing.T) {
	// valid name
	err := checkName("marcus")
	if err != nil {
		t.Fatalf(`checkName("marco") results in error`)
	}

	// name with space
	name := "marcus aurelius "
	name = stripSpace(name)
	err = checkName(name)
	if err != nil {
		t.Fatalf(`A name containing spaces results in error`)
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

func TestRemoveVowels(t *testing.T) {
	s := "AEIOUX"
	s = removeVowels(s)
	if s != "X" {
		t.Fatalf(`Remove vowels function failure`)
	}
}
