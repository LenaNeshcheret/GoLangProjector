package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	data, err := os.ReadFile("1689007675141_numbers.txt")
	if err != nil {
		log.Fatal(err)
	}

	phoneWithDigit10 := `\d{10}`
	phoneWithParentheses := `\(\d{3}\)\s*\d{3}[-.\s]?\d{4}`
	phoneWithSpaces := `\d{3}\s\d{3}\s\d{4}`
	phoneWithHyphens := `\d{3}-\d{3}-\d{4}`
	phoneWithAll := `(\(\d{3}\)\s*|\d{3}[-.\s]?)?\d{3}[-.\s]?\d{4}`

	printPhoneNumbers(data, phoneWithDigit10, "10-digit phone numbers")
	printPhoneNumbers(data, phoneWithParentheses, "Phone numbers with parentheses")
	printPhoneNumbers(data, phoneWithSpaces, "Phone numbers with spaces")
	printPhoneNumbers(data, phoneWithHyphens, "Phone numbers with hyphens")
	printPhoneNumbers(data, phoneWithAll, "All phone number formats")
}

func findPhoneNumbers(data string, pattern string) []string {
	re := regexp.MustCompile(pattern)

	phones := re.FindAllString(data, -1)

	return phones
}

func printPhoneNumbers(data []byte, pattern string, description string) {
	fmt.Printf("Finding %s:\n", description)
	phones := findPhoneNumbers(string(data), pattern)
	for _, phone := range phones {
		fmt.Println(phone)
	}
	fmt.Println()
}
