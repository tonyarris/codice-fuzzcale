package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/common-nighthawk/go-figure"
)

func main() {

	// print title
	title := figure.NewColorFigure("codice FUZZcale", "slant", "red", true)
	title.Print()
	time.Sleep(1000)
	fmt.Println()

	// collect known information from user
	// TODO - implement unknown functionality
	// prompt & store surname
	fmt.Println("Enter surname(s) (Enter for unknown): ")
	reader := bufio.NewReader(os.Stdin)
	surname, _ := reader.ReadString('\n')

	// replace newline
	surname = replaceNewLine(surname)

	// strip spaces
	surname = stripSpace(surname)

	// set logs
	log.SetPrefix("surname: ")
	log.SetFlags(0)

	// check validity
	err := checkName(surname)
	if err != nil {
		log.Fatal(err)
	}
	surname = strings.ToUpper(surname)

	// prompt & store firstname
	fmt.Println("Enter firstname(s) (Enter for unknown): ")
	firstname, _ := reader.ReadString('\n')

	// replace newline
	firstname = replaceNewLine(firstname)
	// strip spaces
	firstname = stripSpace(firstname)

	// set logs
	log.SetPrefix("firstname: ")
	log.SetFlags(0)

	// check validity
	err = checkName(firstname)
	if err != nil {
		log.Fatal(err)
	}
	firstname = strings.ToUpper(surname)

	// prompt & store sex
	fmt.Println("Enter sex (M/F/Enter for unknown): ")
	var sex string
	sex, _ = reader.ReadString('\n')
	sex = strings.ToUpper(sex)

	// set logs
	log.SetPrefix("sex: ")
	log.SetFlags(0)

	// validate sex
	var sexCheck []rune
	sexCheck = []rune(sex)
	err = checkSex(sexCheck)
	if err != nil {
		log.Fatal(err)
	}

	// prompt & store DOB
	const (
		layoutISO = "2006-01-02"
	)
	fmt.Println("Enter date of birth (yyyy-mm-dd): ")
	var dob string
	fmt.Scanln(&dob)
	t, _ := time.Parse(layoutISO, dob)

	// prompt & store comune
	fmt.Println("Enter comune of birth: ")
	var comune string
	fmt.Scanln(&comune)
	comune = strings.ToUpper(comune)

	// Construct codice fiscale

	// surname triplet
	// remove vowels from surname
	surname = removeVowels(surname)

	if len(surname) >= 3 {
		surname = surname[0:3]
	}
	// TODO - account for <3 consonants & <3 letters

	// name triplet
	// remove vowels from name
	firstname = removeVowels(firstname)

	// if > 3 consonants in firstname, skip the second
	var nameTrip []rune = []rune(firstname)
	if len(firstname) > 3 {
		nameTrip = delChar(nameTrip, 1)
	}
	firstname = string(nameTrip)

	if len(firstname) >= 3 {
		firstname = firstname[0:3]
	}

	// birth year
	var birthYear int
	birthYear = t.Year()
	birthYear = birthYear % 100

	// birth month map
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

	// comune code
	// read comune names
	content, err := ioutil.ReadFile("./comune_codes/final_codes/comune_names.txt")
	if err != nil {
		log.Fatal(err)
	}
	cNames := splitString(content)

	// read comune codes
	content2, err2 := ioutil.ReadFile("./comune_codes/final_codes/comune_codes.txt")
	if err2 != nil {
		log.Fatal(err)
	}
	cCodes := splitString(content2)

	// TODO - add foreign comune codes

	// create comune map
	comuneMap := make(map[string]string)
	for i := range cNames {

		comuneMap[(cNames[i])] = cCodes[i]
	}

	comuneCode := comuneMap[comune]

	// calculate check character
	var cf string = surname + firstname + strconv.Itoa(birthYear) + mCode + strconv.Itoa(day) + comuneCode

	// split into evens & evens
	runeCF := []rune(cf)
	var odd string
	var even string

	for i := 0; i < len(runeCF); i = i + 2 {
		odd = odd + string(runeCF[i])
	}
	for i := 1; i < len(runeCF); i = i + 2 {
		even = even + string(runeCF[i])
	}

	// odd check map
	oddMap := make(map[string]int)
	oddMap["0"] = 1
	oddMap["1"] = 0
	oddMap["2"] = 5
	oddMap["3"] = 7
	oddMap["4"] = 9
	oddMap["5"] = 13
	oddMap["6"] = 15
	oddMap["7"] = 17
	oddMap["8"] = 19
	oddMap["9"] = 21
	oddMap["A"] = 1
	oddMap["B"] = 0
	oddMap["C"] = 5
	oddMap["D"] = 7
	oddMap["E"] = 9
	oddMap["F"] = 13
	oddMap["G"] = 15
	oddMap["H"] = 17
	oddMap["I"] = 19
	oddMap["J"] = 21
	oddMap["K"] = 2
	oddMap["L"] = 4
	oddMap["M"] = 18
	oddMap["N"] = 20
	oddMap["O"] = 11
	oddMap["P"] = 3
	oddMap["Q"] = 6
	oddMap["R"] = 8
	oddMap["S"] = 12
	oddMap["T"] = 14
	oddMap["U"] = 16
	oddMap["V"] = 10
	oddMap["W"] = 22
	oddMap["X"] = 25
	oddMap["Y"] = 24
	oddMap["Z"] = 23

	// even check map
	evenMap := make(map[string]int)
	evenMap["0"] = 0
	evenMap["1"] = 1
	evenMap["2"] = 2
	evenMap["3"] = 3
	evenMap["4"] = 4
	evenMap["5"] = 5
	evenMap["6"] = 6
	evenMap["7"] = 7
	evenMap["8"] = 8
	evenMap["9"] = 9
	evenMap["A"] = 0
	evenMap["B"] = 1
	evenMap["C"] = 2
	evenMap["D"] = 3
	evenMap["E"] = 4
	evenMap["F"] = 5
	evenMap["G"] = 6
	evenMap["H"] = 7
	evenMap["I"] = 8
	evenMap["J"] = 9
	evenMap["K"] = 10
	evenMap["L"] = 11
	evenMap["M"] = 12
	evenMap["N"] = 13
	evenMap["O"] = 14
	evenMap["P"] = 15
	evenMap["Q"] = 16
	evenMap["R"] = 17
	evenMap["S"] = 18
	evenMap["T"] = 19
	evenMap["U"] = 20
	evenMap["V"] = 21
	evenMap["W"] = 22
	evenMap["X"] = 23
	evenMap["Y"] = 24
	evenMap["Z"] = 25

	// calculate intermediate values
	runeOdd := []rune(odd)
	runeEven := []rune(even)
	var oddValue int
	var evenValue int
	for i := 0; i < len(runeOdd); i++ {
		oddValue = oddValue + oddMap[string(runeOdd[i])]
	}
	for i := 0; i < len(runeEven); i++ {
		evenValue = evenValue + evenMap[string(runeEven[i])]
	}
	combinedValue := (oddValue + evenValue) % 26

	// remainder map
	rem := make(map[int]string)
	rem[0] = "A"
	rem[1] = "B"
	rem[2] = "C"
	rem[3] = "D"
	rem[4] = "E"
	rem[5] = "F"
	rem[6] = "G"
	rem[7] = "H"
	rem[8] = "I"
	rem[9] = "J"
	rem[10] = "K"
	rem[11] = "L"
	rem[12] = "M"
	rem[13] = "N"
	rem[14] = "O"
	rem[15] = "P"
	rem[16] = "Q"
	rem[17] = "R"
	rem[18] = "S"
	rem[19] = "T"
	rem[20] = "U"
	rem[21] = "V"
	rem[22] = "W"
	rem[23] = "X"
	rem[24] = "Y"
	rem[25] = "Z"

	check := rem[combinedValue]

	// print concatenated CF
	cf = cf + check
	fmt.Print(cf)

}

func delChar(s []rune, index int) []rune {
	return append(s[0:index], s[index+1:]...)
}

func checkName(s string) error {
	// sanify input
	// check contains only letters & spaces
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return errors.New("INVALID INPUT")
		}
	}
	return nil
}

func stripSpace(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func removeVowels(s string) string {
	for _, c := range []string{"A", "E", "I", "O", "U"} {

		s = strings.ReplaceAll(s, c, "")
	}
	return s
}

func replaceNewLine(s string) string {
	if runtime.GOOS == "windows" {
		// for Windows compatibility
		s = strings.Replace(s, "\r\n", "", -1)
	} else {
		s = strings.Replace(s, "\n", "", -1)
	}
	return s
}

func splitString(b []byte) []string {
	var splitList []string
	if runtime.GOOS == "windows" {
		// for Windows compatibility
		splitList = strings.Split(string(b), "\r\n")
	} else {
		splitList = strings.Split(string(b), "\n")
	}
	return splitList
}

func checkSex(s []rune) error {
	if s[0] != 'M' {
		if s[0] != 'F' {
			if len(s) > 1 {
				return errors.New("INVALID SEX")
			}
		}
	}

	return nil
}
