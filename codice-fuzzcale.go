package main

import (
	"fmt"
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

	// prompt & store firstname
	fmt.Println("Enter firstname(s) (Enter for unknown): ")
	var firstname string
	fmt.Scanln(&firstname)

	// prompt & store sex
	fmt.Println("Enter sex (M/F/Enter for unknown): ")
	var sex string
	fmt.Scanln(&sex)

	// prompt & store DOB
	const (
		layoutISO = "2006-01-02"
	)
	fmt.Println("Enter date of birth (yyyy-mm-dd): ")
	var dob string
	fmt.Scanln(&dob)
	t, _ := time.Parse(layoutISO, dob)
	fmt.Print(t)

	// prompt & store comune
	fmt.Println("Enter comune of birth: ")
	var comune string
	fmt.Scanln(&comune)

	//TODO - calculate check character

}
