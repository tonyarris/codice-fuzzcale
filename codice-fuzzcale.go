package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
)

func main() {

	// print title
	title := figure.NewColorFigure("codice FUZZcale", "slant", "red", true)
	title.Print()
	time.Sleep(1)
	fmt.Println()

	// collect known information from user
	// TODO - implement unknown functionality
	// prompt & store surname
	fmt.Println("Enter surname(s) (Enter for unknown): ")
	var surname string
	fmt.Scanln(&surname)
	surname = strings.ToUpper(surname)

	// prompt & store firstname
	fmt.Println("Enter firstname(s) (Enter for unknown): ")
	var firstname string
	fmt.Scanln(&firstname)
	firstname = strings.ToUpper(firstname)

	// prompt & store sex
	fmt.Println("Enter sex (M/F/Enter for unknown): ")
	var sex string
	sex = strings.ToUpper(sex)
	fmt.Scanln(&sex)

	// prompt & store DOB
	const (
		layoutISO = "2006-01-02"
	)
	fmt.Println("Enter date of birth (yyyy-mm-dd): ")
	var dob string
	fmt.Scanln(&dob)
	t, _ := time.Parse(layoutISO, dob)
	//fmt.Print(t)

	// prompt & store comune
	fmt.Println("Enter comune of birth: ")
	var comune string
	fmt.Scanln(&comune)

	// Construct codice fiscale

	// surname triplet
	// remove vowels from surname
	for _, c := range []string{"A", "E", "I", "O", "U"} {

		surname = strings.ReplaceAll(surname, c, "")
	}
	surname = surname[0:3]
	// TODO - account for <3 consonants & <3 letters

	// name triplet
	// remove vowels from name
	for _, c := range []string{"A", "E", "I", "O", "U"} {

		firstname = strings.ReplaceAll(firstname, c, "")
	}

	// if > 3 consonants in firstname, skip the second
	var nameTrip []rune = []rune(firstname)
	if len(firstname) > 3 {
		nameTrip = delChar(nameTrip, 1)
	}
	firstname = string(nameTrip)
	firstname = firstname[0:3]

	// birth year
	var birthYear int
	birthYear = t.Year()
	birthYear = birthYear % 100

	// birth month dict
	m := make(map[string]string)
	m["January"] = "A"
	m["February"] = "B"
	m["March"] = "C"
	m["April"] = "D"
	m["May"] = "E"
	m["June"] = "H"
	m["July"] = "L"
	m["August"] = "M"
	m["September"] = "P"
	m["October"] = "R"
	m["November"] = "S"
	m["December"] = "T"

	// get month code
	mCode := m[t.Month().String()]

	// day of birth
	// day counter
	var dayCount int
	dayCount = 0

	if sex == "F" {
		dayCount = 40
	}
	// actual day of birth, plus 40 for F
	var day int = t.Day() + dayCount

	//TODO - calculate check character

	// print concatenated CF
	fmt.Print(surname, firstname, birthYear, mCode, day)

}

func delChar(s []rune, index int) []rune {
	return append(s[0:index], s[index+1:]...)
}
