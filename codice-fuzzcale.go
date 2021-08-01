package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/common-nighthawk/go-figure"
)

// birth month map
var m = map[string]string{
	"January":   "A",
	"February":  "B",
	"March":     "C",
	"April":     "D",
	"May":       "E",
	"June":      "H",
	"July":      "L",
	"August":    "M",
	"September": "P",
	"October":   "R",
	"November":  "S",
	"December":  "T",
}

// odd check map
var oddMap = map[string]int{
	"0": 1,
	"1": 0,
	"2": 5,
	"3": 7,
	"4": 9,
	"5": 13,
	"6": 15,
	"7": 17,
	"8": 19,
	"9": 21,
	"A": 1,
	"B": 0,
	"C": 5,
	"D": 7,
	"E": 9,
	"F": 13,
	"G": 15,
	"H": 17,
	"I": 19,
	"J": 21,
	"K": 2,
	"L": 4,
	"M": 18,
	"N": 20,
	"O": 11,
	"P": 3,
	"Q": 6,
	"R": 8,
	"S": 12,
	"T": 14,
	"U": 16,
	"V": 10,
	"W": 22,
	"X": 25,
	"Y": 24,
	"Z": 23,
}

// even check map
var evenMap = map[string]int{
	"0": 0,
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"A": 0,
	"B": 1,
	"C": 2,
	"D": 3,
	"E": 4,
	"F": 5,
	"G": 6,
	"H": 7,
	"I": 8,
	"J": 9,
	"K": 10,
	"L": 11,
	"M": 12,
	"N": 13,
	"O": 14,
	"P": 15,
	"Q": 16,
	"R": 17,
	"S": 18,
	"T": 19,
	"U": 20,
	"V": 21,
	"W": 22,
	"X": 23,
	"Y": 24,
	"Z": 25,
}

// remainder map
var rem = map[int]string{
	0:  "A",
	1:  "B",
	2:  "C",
	3:  "D",
	4:  "E",
	5:  "F",
	6:  "G",
	7:  "H",
	8:  "I",
	9:  "J",
	10: "K",
	11: "L",
	12: "M",
	13: "N",
	14: "O",
	15: "P",
	16: "Q",
	17: "R",
	18: "S",
	19: "T",
	20: "U",
	21: "V",
	22: "W",
	23: "X",
	24: "Y",
	25: "Z",
}

var comuneMap = createComuneMap()

// unknown variable detection bools
var fuzzSurname, fuzzFirstname, fuzzSex, fuzzDob, fuzzComune bool

func main() {

	// print title
	title := figure.NewColorFigure("codice FUZZcale", "slant", "red", true)
	title.Print()
	fmt.Println()

	// collect known information from user
	fmt.Println("Just hit ENTER for unknown values")
	// prompt & store surname
	fmt.Println("Enter surname(s): ")
	reader := bufio.NewReader(os.Stdin)
	surname, _ := reader.ReadString('\n')

	// replace newline
	surname = replaceNewLine(surname)

	// strip spaces
	surname = stripSpace(surname)

	// detect unknown
	if len(surname) == 0 {
		fuzzSurname = true
	}

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
	fmt.Println("Enter firstname(s):")
	firstname, _ := reader.ReadString('\n')

	// replace newline
	firstname = replaceNewLine(firstname)
	// strip spaces
	firstname = stripSpace(firstname)

	// detect unknown
	if len(firstname) == 0 {
		fuzzFirstname = true
		fmt.Print(fuzzFirstname)
	}

	// set logs
	log.SetPrefix("firstname:")
	log.SetFlags(0)

	// check validity
	err = checkName(firstname)
	if err != nil {
		log.Fatal(err)
	}
	firstname = strings.ToUpper(firstname)

	// prompt & store sex
	fmt.Println("Enter sex (M/F):")
	var sex string
	sex, _ = reader.ReadString('\n')
	sex = strings.ToUpper(sex)
	sex = replaceNewLine(sex)

	// detect unknown
	if len(sex) == 0 {
		fuzzSex = true
		fmt.Print(fuzzSex)
	}

	// set logs
	log.SetPrefix("sex:")
	log.SetFlags(0)

	// validate sex
	err = checkSex(sex)
	if err != nil {
		log.Fatal(err)
	}

	// prompt & store DOB
	const (
		layoutISO = "2006-01-02"
	)
	fmt.Println("Enter date of birth (yyyy-mm-dd):")
	var dob string
	dob, _ = reader.ReadString('\n')
	dob = replaceNewLine(dob)
	t, _ := time.Parse(layoutISO, dob)

	// detect unknown
	if len(dob) == 0 {
		fuzzDob = true
		fmt.Print(fuzzDob)
	}

	// set logs
	log.SetPrefix("date of birth:")
	log.SetFlags(0)

	// TODO verify birth date format
	//err = checkDate(dob)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// prompt & store comune
	fmt.Println("Enter comune of birth:")
	var comune string
	comune, _ = reader.ReadString('\n')
	comune = replaceNewLine(comune)
	comune = strings.ToUpper(comune)

	// detect unknown
	if len(comune) == 0 {
		fuzzComune = true
		fmt.Print(fuzzComune)
	}

	// TODO prompt for output file location

	// Construct codice fiscale from known values

	if !fuzzSurname {
		// surname triplet
		// extract vowels from surname
		s_vowels := extractVowels(surname)
		// remove vowels from surname
		surname = removeVowels(surname)

		surname = constructTriplet(surname, s_vowels)
	}

	if !fuzzFirstname {
		// name triplet
		// extract vowels from firstname
		f_vowels := extractVowels(firstname)
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

		firstname = constructTriplet(firstname, f_vowels)
	}

	var mCode string
	var birthYear, day int
	if !fuzzDob {
		// birth year
		birthYear = t.Year()
		birthYear = birthYear % 100

		// get month code
		mCode = m[t.Month().String()]

		// day of birth
		// day counter
		var dayCount int
		dayCount = 0

		if sex == "F" {
			dayCount = 40
		}
		// actual day of birth, plus 40 for F
		day = t.Day() + dayCount
	}

	var comuneCode string
	if !fuzzComune {
		// assign comune code
		comuneCode = comuneMap[comune]
	}

	// calculate single cf with all known values
	var cf, check string
	if !fuzzSurname && !fuzzFirstname && !fuzzSex && !fuzzDob && !fuzzComune {
		// construct cf minus check
		cf = constructCFNoCheck(surname, firstname, birthYear, mCode, day, comuneCode)
		// calculate check character
		check = calculateCheck(cf)
		// print concatenated CF
		cf = replaceNewLine(cf + check)
		fmt.Print(cf)
	} else if fuzzSurname { // if surname unknown
		c := make(chan [3]string)
		go fuzzAlphabet(c)
		for sur := range c {
			surname = sur[0] + sur[1] + sur[2]
			// construct cf minus check
			cf = constructCFNoCheck(surname, firstname, birthYear, mCode, day, comuneCode)
			// calculate check character
			check = calculateCheck(cf)
			// print concatenated CF
			cf = replaceNewLine(cf + check)
			fmt.Println(cf)
		}
	}

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

func extractVowels(s string) string {
	for _, c := range []string{"B", "C", "D", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "V", "W", "X", "Y", "Z"} {

		s = strings.ReplaceAll(s, c, "")
	}
	return s
}

func constructTriplet(s, v string) string {
	if len(s) >= 3 {
		s = s[0:3]
	} else {
		s = s + v
		// if < 3 letters total, fill with X
		if len(s) > 0 && len(s) < 3 {
			s = s + "XX"
			s = s[0:3]
		}
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

func checkSex(s string) error {
	if s == "M" || s == "F" || s == "" {
		return nil
	} else {
		return errors.New("INVALID SEX")
	}
}

// TODO - fix check
func checkDate(s string) error {
	fmt.Print(s)
	reg, _ := regexp.Compile(`/^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$/`)
	matched := reg.MatchString(s)
	if !matched {
		return errors.New("INVALID DOB")
	} else {
		return nil
	}
}

func createComuneMap() map[string]string {
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

	// create comune map
	comuneMap := make(map[string]string)
	for i := range cNames {

		comuneMap[(cNames[i])] = cCodes[i]
	}

	return comuneMap
}

func fuzzAlphabet(c chan [3]string) {
	var triplet [3]string

	for ch := 'A'; ch <= 'Z'; ch++ {
		triplet[0] = string(ch)
		for ch := 'A'; ch <= 'Z'; ch++ {
			triplet[1] = string(ch)
			for ch := 'A'; ch <= 'Z'; ch++ {
				triplet[2] = string(ch)
				c <- triplet
			}
		}
	}
	close(c)
}

func calculateCheck(s string) string {
	// split into odds & evens
	runeCF := []rune(s)
	var odd string
	var even string

	for i := 0; i < len(runeCF); i = i + 2 {
		odd = odd + string(runeCF[i])
	}
	for i := 1; i < len(runeCF); i = i + 2 {
		even = even + string(runeCF[i])
	}

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

	check := rem[combinedValue]
	return check
}

func constructCFNoCheck(surname string, firstname string, birthYear int, mCode string, day int, comuneCode string) string {
	// construct cf minus check
	cf := surname + firstname + strconv.Itoa(birthYear) + mCode + fmt.Sprintf("%02d", day) + comuneCode
	return cf
}
