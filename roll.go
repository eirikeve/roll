package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Seed for rng
	rand.Seed(int64(time.Now().Nanosecond()))
	// Setup flag
	helpFlagPtr := flag.Bool("h", false, "Help flag")
	flag.Parse()
	if *helpFlagPtr {
		printHelpMessage()
		return
	}
	// Read cmd line arguments
	args := strings.Join(os.Args[1:], " ") // Roll arguments

	if !isArgumentAcceptable(&args) {
		fmt.Println("Unacceptable argument ", args)
		fmt.Println("Use -h for help")
		return
	}
	integerStrings, dicethrowStrings := formatInput(&args)
	printInputConfirmation(integerStrings, dicethrowStrings)
	fmt.Println()
	integers, dicethrows := getResults(integerStrings, dicethrowStrings)
	printThrowResults(dicethrowStrings, dicethrows)
	sumIntegers, sumThrows := sumResults(integers, dicethrows)
	fmt.Println("Const:\t", sumIntegers)
	fmt.Println("Sum: \t", sumIntegers+sumThrows)
}

/*
Returns true if cmd line arguments are acceptable
(only includes dice rolls), false else.
*/
func isArgumentAcceptable(a *string) bool {
	if isArgTermsSeparated(a) && isArgAcceptableRunes(a) {
		return true
	}
	fmt.Println("isArgTermsSeparated: ", isArgTermsSeparated(a))
	fmt.Println("isArgAcceptableRunes: ", isArgAcceptableRunes(a))
	return false
}

/*
Returns true if all runes (chars) are acceptable chars, false else.
*/
func isArgAcceptableRunes(a *string) bool {
	acceptable := []rune{' ', 'd', '+', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	for _, rune := range *a {
		var match bool = false
		for _, b := range acceptable {
			if rune == b {
				match = true
				break
			}
		}
		if !(match) {
			return false
		}
	}
	return true
}

/*
Returns true if all dice terms are separated, false else
Uses regexp to do this.
The regexp 3d33d3 and 3dd3, but not 3d3 3d3 or 3d3+3d3, for instance
*/
func isArgTermsSeparated(a *string) bool {
	r, _ := regexp.Compile("d[0-9]+d|dd")
	match := r.MatchString(*a)
	return (!match)
}

/*
Throws virtual dice
*/
func throw(numberOfDice int, sides int) int {
	if numberOfDice == 0 || sides == 0 {
		return 0
	}
	var total int = 0
	for i := 0; i < numberOfDice; i++ {
		total += ((rand.Int() % sides) + 1)
	}
	return total
}

func formatInput(a *string) ([]string, []string) {
	// Matches integers not part of a dice throw
	intRegExp, _ := regexp.Compile("[^d0-9][0-9]+[^d0-9]|-?[^d0-9][0-9]+$")
	// Matches dice throws
	diceRegExp, _ := regexp.Compile("^[0-9]+d[0-9]+|[^0-9]+[0-9]+d[0-9]+")
	integerStrings := intRegExp.FindAllString(*a, -1)
	dicethrowStrings := diceRegExp.FindAllString(*a, -1)
	for i, _ := range integerStrings {
		integerStrings[i] = strings.Replace(integerStrings[i], " ", "", -1)
		integerStrings[i] = strings.Replace(integerStrings[i], "+", "", -1)
	}
	for i, _ := range dicethrowStrings {
		dicethrowStrings[i] = strings.Replace(dicethrowStrings[i], " ", "", -1)
		dicethrowStrings[i] = strings.Replace(dicethrowStrings[i], "+", "", -1)
	}

	return integerStrings, dicethrowStrings
}

func getThrowFromString(dicethrow string) int {
	var isNegative bool = false
	var delimiterIndex int = 0
	// Check if negative & if so, remove sign from string
	if dicethrow[0] == '-' {
		isNegative = true
		dicethrow = dicethrow[1:]
	}
	// Find where the 'd' is
	for i, _ := range dicethrow {
		if dicethrow[i] == 'd' {
			delimiterIndex = i
			break
		}
	}
	numberOfDice, _ := strconv.Atoi(dicethrow[0:delimiterIndex])
	sides, _ := strconv.Atoi(dicethrow[delimiterIndex+1:])
	result := throw(numberOfDice, sides)
	if isNegative {
		result = -result
	}
	return result
}

func getResults(integerStrings []string, dicethrowStrings []string) ([]int, []int) {
	var integers []int
	var dicethrows []int

	for _, s := range integerStrings {
		n, _ := strconv.Atoi(s)
		integers = append(integers, n)
	}
	for _, s := range dicethrowStrings {
		dicethrows = append(dicethrows, getThrowFromString(s))
	}
	return integers, dicethrows

}

func sumResults(integers []int, dicethrows []int) (int, int) {
	var sumIntegers int = 0
	var sumThrows int = 0
	for _, val := range integers {
		sumIntegers += val
	}
	for _, val := range dicethrows {
		sumThrows += val
	}
	return sumIntegers, sumThrows
}

func printInputConfirmation(integers []string, dice []string) {
	print("Rolling: ")
	for _, v := range dice {
		fmt.Print(v)
		fmt.Print(" ")
	}
	for i, v := range integers {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}

	fmt.Println()
}

func printThrowResults(dicethrowStrings []string, dicethrows []int) {
	fmt.Print("Throws:")
	for i, s := range dicethrowStrings {
		fmt.Println("\t", s, "\t->", dicethrows[i])
	}
}

/*
Implementation of self-destruct functionality.
*/
func printHelpMessage() {
	fmt.Println("/--- roll command line utility ---/")
	fmt.Println("roll implements command line dice rolling.")
	fmt.Println("Input should consist of any number of terms.")
	fmt.Println("Terms can be either a constant or dice throw, and they are separated by + or -.")
	fmt.Println("If there is no sign before a term, it is assumed to be positive.")
	fmt.Println("Dice are written as adb, where we throw a dice with b sides. (Use lowercase d!)")
	fmt.Println("Constants are any number")
	fmt.Println("Example input: roll 3d20 + 5 - 1d4")
}
