package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"unicode"
	"unicode/utf8"
)

func main() {

	data, err := os.ReadFile("1689007676028_text.txt")
	if err != nil {
		log.Fatal(err)
	}
	text := string(data)

	// Regex pattern to match words starting and ending with specified Ukrainian vowels
	pattern := `[аеєиіїоуюяАЕЄИІЇОУЮЯ][а-яА-Я]*[аеєиіїоуюяАЕЄИІЇОУЮЯ]`

	re := regexp.MustCompile(pattern)

	// Find all matches
	matches := re.FindAllStringIndex(text, -1)

	var filteredMatches []string
	for _, match := range matches {
		// Check if the match is at the start of the text or preceded by a non-letter character
		if (match[0] == 0 || !isLetterBeforeIndex(text, match[0])) && (match[1] == len(text) || !isLetterAfterIndex(text, match[1])) {
			filteredMatches = append(filteredMatches, text[match[0]:match[1]])
		}
	}
	fmt.Println(filteredMatches)
}

// Helper function to check if the character before the given index is a letter
func isLetterBeforeIndex(text string, index int) bool {
	if index == 0 {
		return false
	}
	runeBefore, _ := utf8.DecodeLastRuneInString(text[:index])
	return unicode.IsLetter(runeBefore)
}

// Helper function to check if the character after the given index is a letter
func isLetterAfterIndex(text string, index int) bool {
	if index >= len(text) {
		return false
	}
	runeAfter, _ := utf8.DecodeRuneInString(text[index:])
	return unicode.IsLetter(runeAfter)
}
