// roll implements command line dice rolling.
//
// It accepts the standard notation, i.e. 10d20 + 5, or 2d7 + 1d10 - 1 etc.,
// and rejects arguments that contain other letters or have unseparated throws (like 3d33d3)
// The calculation is done by separating the dice throws and constants (using regexp),
// Computing the result of each throw using pRNG,
// and summing the results.
//
// Written by Eirik Vesterkjær, dec 2017

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
	// Check for -h flag, which prompts printHelpMessage()
	helpFlagPtr := flag.Bool("h", false, "Help flag")
	flag.Parse()
	if *helpFlagPtr {
		printHelpMessage()
		return
	}
	// Read cmd line arguments and concatenate
	args := strings.Join(os.Args[1:], " ")

	// Check for validity of arguments, print error msg if not acceptable
	if !isArgumentAcceptable(&args) {
		fmt.Println("Unacceptable argument ", args, ":")
		if !isArgTermsSeparated(&args) {
			fmt.Println("- Lacking separation between dice throws.")
		}
		if !isArgAcceptableRunes(&args) {
			fmt.Println("- Non-accepted characters in input.")
		}
		fmt.Println("Use the -h flag for help")
		return
	}
	// Separates dice throws and constants, while preserving negative signs
	integerStrings, dicethrowStrings := formatInput(&args)
	printInputConfirmation(integerStrings, dicethrowStrings)
	// Throws dice, converts constants to int.
	integers, dicethrows := getResults(integerStrings, dicethrowStrings)
	printThrowResults(dicethrowStrings, dicethrows)
	// Separate sums in order to display the sum of the constants (integers)
	sumIntegers, sumThrows := sumResults(integers, dicethrows)
	fmt.Println("Const:\t", sumIntegers)
	fmt.Println("Sum: \t", sumIntegers+sumThrows)
}

/*
Checks if cmd line arguments are acceptable (i.e. a dice throw cmd)
@arg a: Ptr to concatenated cmd line argument
@return: true if cmd line arguments are acceptable, else false
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
Checks if all characters in cmd line argument are either: 0-9, whitespace, d, +, -
@arg a: Ptr to concatenated cmd line argument
@return: true if all chars are acceptable, else false
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
Checks is all dice terms are separated, using regexp.
These are not separated: 	"3d33d3" and "3dd3"
These are separated: 		"3d3 3d3" and "3d3+3d3"
@arg a: Ptr to concatenated cmd line argument
@return: true if all throws are separated, else false
*/
func isArgTermsSeparated(a *string) bool {
	r, _ := regexp.Compile("d[0-9]+d|dd")
	match := r.MatchString(*a)
	return (!match)
}

/*
Throws virtual dice using pRNG
@arg numberOfDice: number of times the die is thrown
@arg sides: Maximum value of the die. Assumes die values [1, sides]
@return: sum of the throws
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

/*
Formats the input cmd line argument. Does this by matching terms using regexp.
Creates a list of dice throws, and a list of constants.
@arg a: Ptr to concatenated cmd line argument
@return: list of constants, list of dice throws
*/
func formatInput(a *string) ([]string, []string) {
	// Insert space before each operator to make regexp read correctly
	*a = strings.Replace(*a, "+", " +", -1)
	*a = strings.Replace(*a, "-", " -", -1)
	// Matches integers not part of a dice throw
	intRegExp, _ := regexp.Compile("-?[^d0-9][0-9]+[^d0-9]|-?[^d0-9][0-9]+$")
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

/*
Given a dice throw input of type string, determines the
number of dice to throw, and the number of sides on the dice, and
throws the dice using throw(..)
@arg dicethrow: String representing a dice throw. E.g. "3d5" or "-6d10"
@return: result of the dice throw
*/
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

/*
Takes the separated constants and dice throws string lists,
and creates two new lists filled with the integer values of the constants,
and the integer results of the dice throws.
@arg integerStrings: Array of integers on string form
@arg dicethrowStrings: Array of dice throws on string form
@return list of integers as int, and list of dice throw result as int
*/
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

/*
Sums two lists and returns their individual sums
@arg integers: list of integers
@arg dicethrows: another list of integers
@return int sum of intergers, int sum of dicethrows
*/
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

/*
Prints the interpreted cmd line argument.
@arg integerStrings: List of constant integers on string form
@arg dicethrowStrings: List of dice throw string representations
*/
func printInputConfirmation(integerStrings []string, dicethrowStrings []string) {
	print("Rolling: ")
	for _, v := range dicethrowStrings {
		fmt.Print(v)
		fmt.Print(" ")
	}
	for i, v := range integerStrings {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}

	fmt.Println()
}

/*
Prints the dice throws string representations and their result side by side
@arg dicethrowStrings: List of dice throw string representations
@arg dicethrows: List of dice throw results
*/
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
	fmt.Println("Negative signs before a term can be used!")
	fmt.Println("Example input: roll 3d20 + 5 - 1d4 - 6")

}
