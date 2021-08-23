package main

import (
	"fmt"
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
		t.Fatalf(`removeVowels() failure`)
	}
}

func TestCheckSex(t *testing.T) {
	// valid input
	err := checkSex("M")
	if err != nil {
		t.Fatal("checkSex() failure")
	}

	err = checkSex("F")
	if err != nil {
		t.Fatal("checkSex() failure")
	}

	err = checkSex("")
	if err != nil {
		t.Fatal("checkSex() failure")
	}

	// invalid input
	err = checkSex("X")
	if err == nil {
		t.Fatal("checkSex() failure - allows invalid input")
	}

	err = checkSex("M1")
	if err == nil {
		t.Fatal("checkSex() failure - allows invalid input")
	}
}

func TestFuzzAlphabet(t *testing.T) {
	c2 := make(chan [3]string)
	go fuzzAlphabet(c2)
}

func TestCheckAges(t *testing.T) {

	// valid input
	err := checkAges(2, 1)
	if err != nil {
		t.Fatal("checkAges() failure - allows invalid input")
	}

	// valid input, ages equal
	err = checkAges(2, 2)
	if err != nil {
		t.Fatal("checkAges() failure - allows invalid input")
	}

	// max lower than min
	err = checkAges(1, 2)
	if err == nil {
		t.Fatal("checkAges() failure - allows invalid input")
	}
}

func TestExtractVowels(t *testing.T) {
	name := "AAAAAAAAAAAAAAAAA"
	name = extractVowels(name)
	if len(name) != 3 {
		t.Fatal("extractVowels returns string of length >3")
	}

	// TODO: FIX UNEXPECTED BEHAVIOUR
	name2 := "ARIHRIR"
	name2 = extractVowels(name2)
	if len(name2) != 3 {
		fmt.Println(len(name2))
		fmt.Println(name2)
		t.Fatal("extractVowels returns vowel string of incorrect length")
	}
}
