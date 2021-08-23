package main

import (
	"bufio"
	"errors"
	"flag"
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

// global CF element vars
var surname, firstname, sex, dob, comuneCode, mCode string
var birthYear, day int
var f *os.File

// define flags
var surnamePtr = flag.String("s", "", "Surname")
var namePtr = flag.String("n", "", "Name")
var sexPtr = flag.String("sex", "", "Sex")
var dobPtr = flag.String("d", "", "Date of birth in the format yyyy-mm-dd")
var comunePtr = flag.String("c", "", "Comune of birth")
var minPtr = flag.Int("min", 0, "Minimum age")
var maxPtr = flag.Int("max", 0, "Maximum age")
var pathPtr = flag.String("p", "", "Output path")

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

var comuneMap, cNames = createComuneMap()

// unknown variable detection bools
var fuzzSurname, fuzzFirstname, fuzzSex, fuzzDob, fuzzComune, maxAge, minAge, writeOut, comuneExist bool
var maxAgeInYears, minAgeInYears int

func main() {

	// parse command line flags
	flag.Parse()

	// print title
	title := figure.NewColorFigure("codice FUZZcale", "slant", "red", true)
	title.Print()

	// prompt if not all flags provided
	if flag.NFlag() != 5 {
		// collect known information from user
		fmt.Println("Just hit ENTER for unknown values")
	}

	// define reader
	reader := bufio.NewReader(os.Stdin)

	if *surnamePtr == "" {
		// prompt & store surname, if not passed as flag
		fmt.Println("Enter surname(s): ")
		surname, _ = reader.ReadString('\n')
	} else {
		surname = *surnamePtr
	}

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
		flag.PrintDefaults()
		log.Fatal(err)
	}
	surname = strings.ToUpper(surname)

	if *namePtr == "" {
		// prompt & store firstname
		fmt.Println("Enter firstname(s):")
		firstname, _ = reader.ReadString('\n')
	} else {
		firstname = *namePtr
	}

	// replace newline
	firstname = replaceNewLine(firstname)
	// strip spaces
	firstname = stripSpace(firstname)

	// detect unknown
	if len(firstname) == 0 {
		fuzzFirstname = true
	}

	// set logs
	log.SetPrefix("firstname:")
	log.SetFlags(0)

	// check validity
	err = checkName(firstname)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}
	firstname = strings.ToUpper(firstname)

	if *sexPtr == "" {
		// prompt & store sex
		fmt.Println("Enter sex (M/F):")
		sex, _ = reader.ReadString('\n')
	} else {
		sex = *sexPtr
	}
	sex = strings.ToUpper(sex)
	sex = replaceNewLine(sex)

	// detect unknown
	if len(sex) == 0 {
		fuzzSex = true
	}

	// set logs
	log.SetPrefix("sex:")
	log.SetFlags(0)

	// validate sex
	err = checkSex(sex)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	// prompt & store DOB
	const (
		layoutISO = "2006-01-02"
	)
	if *minPtr == 0 && *maxPtr == 0 {
		if *dobPtr == "" {
			fmt.Println("Enter date of birth (yyyy-mm-dd):")
			// var dob string
			dob, _ = reader.ReadString('\n')
		} else {
			dob = *dobPtr
		}
		dob = replaceNewLine(dob)
	}

	// detect unknown
	if len(dob) == 0 {
		fuzzDob = true
	}

	// get max/min age if dob unknown
	if fuzzDob {
		if *maxPtr == 0 {
			fmt.Print("Max age: ")
			fmt.Scanf("%d", &maxAgeInYears)
		} else {
			maxAgeInYears = *maxPtr
		}
		if maxAgeInYears > 0 {
			maxAge = true
		}

		if *minPtr == 0 {
			fmt.Print("Min age: ")
			fmt.Scanf("%d", &minAgeInYears)
		} else {
			minAgeInYears = *minPtr
		}
		if minAgeInYears > 0 {
			minAge = true
		}

		if minAge && maxAge {
			//set logs
			log.SetPrefix("Max/min age: ")
			log.SetFlags(0)

			// validate relative ages min/max
			err = checkAges(maxAgeInYears, minAgeInYears)
			if err != nil {
				flag.PrintDefaults()
				log.Fatal(err)
			}
		}
	}

	// set logs
	log.SetPrefix("date of birth:")
	log.SetFlags(0)

	// TODO verify birth date format
	// err = checkDate(dob)
	// if err != nil {
	// 	flag.PrintDefaults()
	// 	log.Fatal(err)
	// }

	// prompt & store comune
	var comune string
	if *comunePtr == "" {
		fmt.Println("Enter comune of birth:")
		comune, _ = reader.ReadString('\n')
	} else {
		comune = *comunePtr
	}
	comune = replaceNewLine(comune)
	comune = strings.ToUpper(comune)

	// detect unknown
	if len(comune) == 0 {
		fuzzComune = true
	}

	// prompt and store output path
	var path string
	if *pathPtr == "" {
		fmt.Println("Output path (/<PATH>/*.txt): ")
		path, _ = reader.ReadString('\n')
	} else {
		path = *pathPtr
	}
	path = replaceNewLine(path)

	// detect outfile
	if len(path) > 0 {
		writeOut = true
	}

	if writeOut {
		// create given file
		f, err = os.Create(path)
		if err != nil {
			flag.PrintDefaults()
			log.Fatal(errors.New("ERROR CREATING OUTFILE"))
		}
		defer f.Close()
	}

	// Construct codice fiscale from known values

	if !fuzzSurname {
		// surname triplet
		// extract vowels from surname
		s_vowels := extractVowels(surname)
		// remove vowels from surname
		surname = removeVowels(surname)

		surname = constructTripletSurname(surname, s_vowels)
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

		firstname = constructTripletFirstname(firstname, f_vowels)
	}

	if !fuzzDob {
		t, _ := time.Parse(layoutISO, dob)

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

	// var comuneCode string
	if !fuzzComune {
		// assign comune code
		comuneCode, comuneExist = comuneMap[comune]

		// set logs
		log.SetPrefix("comune: ")
		log.SetFlags(0)

		// check if comune exists
		if !comuneExist {
			fmt.Print("The comune entered does not exist.\n")
			err = errors.New("COMUNE ERROR")
			log.Fatal(err)
		}
	}

	// calculate single cf with all known values
	if !fuzzSurname && !fuzzFirstname && !fuzzSex && !fuzzDob && !fuzzComune {
		// construct cf
		constructCF(surname, firstname, birthYear, mCode, day, comuneCode, f)
	} else {

		// establish values to be fuzzed
		indicator := generateIndicator()

		// generate CFs
		generateCF(indicator, 0, "")

	}

}

func delChar(s []rune, index int) []rune {
	return append(s[0:index], s[index+1:]...)
}

func checkName(s string) error {
	// sanitise input
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
	sRuneIn := []rune(s)
	var sRuneOut []rune
	for i, v := range sRuneIn {
		switch v {
		case 'A', 'E', 'I', 'O', 'U':
			sRuneOut = append(sRuneOut, sRuneIn[i])
		}
	}

	s = string(sRuneOut)
	if len(s) > 3 {
		s = s[:3]
	}
	return s
}

func constructTripletSurname(s, v string) string {
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
	surname = s
	return s
}

func constructTripletFirstname(s, v string) string {
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
	firstname = s
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

func checkAges(max int, min int) error {
	if max >= min {
		return nil
	} else {
		return errors.New("INVALID AGE RANGE")
	}
}

func createComuneMap() (map[string]string, []string) {
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

	return comuneMap, cNames
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

func fuzzComuneCode(c chan string) {
	for i, s := range cNames {
		c <- comuneMap[s]
		i++
	}
	close(c)
}

func rangeDate(start, end time.Time) func() time.Time {
	y, month, d := end.Date()
	start = time.Date(y, month, d, 0, 0, 0, 0, time.UTC)
	y, month, d = start.Date()
	end = time.Date(y, month, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		date := start
		start = start.AddDate(0, 0, 1)
		birthYear = birthYear % 100
		day = date.Day()
		mCode = m[date.Month().String()]
		return date
	}
}

// calculateCheck() calculates and returns the final check digit of the fiscal code
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

// constructCF() constructs a fiscal code based on complete, known information and fuzzes the sex if necessary
func constructCF(surname string, firstname string, birthYear int, mCode string, day int, comuneCode string, f *os.File) {
	// construct cf minus check
	cf := surname + firstname + strconv.Itoa(birthYear) + mCode + fmt.Sprintf("%02d", day) + comuneCode
	// calculate check character
	check := calculateCheck(cf)
	// print concatenated CF
	cf = replaceNewLine(cf + check)
	fmt.Println(cf)
	f.WriteString(cf + "\n")

	// fuzz sex
	if fuzzSex {
		day = day + 40
		cf = surname + firstname + strconv.Itoa(birthYear) + mCode + fmt.Sprintf("%02d", day) + comuneCode
		check = calculateCheck(cf)
		cf = replaceNewLine(cf)
		fmt.Println(cf + check)
		f.WriteString(cf + check + "\n")
	}
}

// generateIndicator() establishes which values to fuzz and
// returns a list of binary values to be passed to generateCF():
// indicator[0] = surname, indicator[1] = name,
// indicator[2] = dob/sex, indicator[3] = comune
func generateIndicator() []int {
	indicator := []int{0, 0, 0, 0, 0}

	if fuzzSurname {
		indicator[0] = 1
	}
	if fuzzFirstname {
		indicator[1] = 1
	}
	if fuzzDob {
		indicator[2] = 1
	}
	if fuzzSex {
		indicator[3] = 1
	}
	if fuzzComune {
		indicator[4] = 1
	}
	return indicator
}

// generateCF() generates fiscal codes recursively.
// It takes the indicator list, a counter and a string as input.
// With each recursion, it checks whether the element is to be
// fuzzed by consulting the indicator[i] value. If so, it fuzzes
// that value and calls itself, passing the fuzzed value, concatenated
// with the string it received. If not, it simply concatenates the input
// string with the global value provided by the user during runtime.
func generateCF(indicator []int, i int, s string) {
	// fmt.Println("recursion ", i)
	if i == 0 {
		if indicator[i] == 1 {
			// fuzz surname
			c := make(chan [3]string)

			go fuzzAlphabet(c)
			for sur := range c {
				surname = sur[0] + sur[1] + sur[2]
				generateCF(indicator, 1, surname)
			}
		} else {
			generateCF(indicator, 1, surname)
		}
	}
	if i == 1 {
		if indicator[i] == 1 {
			//fuzz firstname
			c2 := make(chan [3]string)
			go fuzzAlphabet(c2)
			for fir := range c2 {
				firstname = fir[0] + fir[1] + fir[2]
				composite := s + firstname
				generateCF(indicator, 2, composite)
			}
		} else {
			generateCF(indicator, 2, s+firstname)
		}
	}
	if i == 2 {
		if indicator[i] == 1 {
			var start, end time.Time
			// set start and end time envelope to match min and max age if entered
			if minAge {
				start = time.Now().AddDate(-minAgeInYears, 0, 0)
			} else {
				start = time.Now()
			}
			if maxAge {
				end = time.Now().AddDate(-maxAgeInYears, 0, 0)
			} else {
				end = time.Now().AddDate(-80, 0, 0)
			}
			for rd := rangeDate(start, end); ; {
				daterange := rd()
				// bottom-out date range
				if daterange.Year() <= start.Year() {
					birthYear := daterange.Year()
					birthYear = birthYear % 100
					day := daterange.Day()
					mCode := m[daterange.Month().String()]
					if !fuzzSex {
						composite := s + strconv.Itoa(birthYear) + mCode + fmt.Sprintf("%02d", day)
						generateCF(indicator, 3, composite)
					} else {
						composite := s + strconv.Itoa(birthYear) + mCode
						generateCF(indicator, 3, composite)
					}
				} else {
					break
				}
			}
		} else {
			if fuzzSex {
				generateCF(indicator, 3, s+strconv.Itoa(birthYear)+mCode)
			} else {
				generateCF(indicator, 3, s+strconv.Itoa(birthYear)+mCode+fmt.Sprintf("%02d", day))
			}
		}
	}
	if i == 3 {
		// TODO - fix fuzz sex
		if fuzzSex {
			composite := s + fmt.Sprintf("%02d", day+40)
			generateCF(indicator, 4, composite)
			composite = s + fmt.Sprintf("%02d", day)
			generateCF(indicator, 4, composite)
		} else {
			generateCF(indicator, 4, s)
		}
	}
	if i == 4 {
		if indicator[i] == 1 {
			c3 := make(chan string)
			go fuzzComuneCode(c3)
			for ccode := range c3 {
				composite := s + ccode
				cf := composite + calculateCheck(composite)
				fmt.Println(cf)
				if writeOut {
					f.WriteString(cf + "\n")
				}
			}
		} else {
			cf := s + comuneCode + calculateCheck(s+comuneCode)
			fmt.Println(cf)
			f.WriteString(cf + "\n")
		}
	}
}
