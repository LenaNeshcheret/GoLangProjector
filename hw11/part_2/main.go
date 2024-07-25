package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {

	data, err := os.ReadFile("1689007676028_text.txt")
	if err != nil {
		log.Fatal(err)
	}
	text := string(data)

	// Regex pattern to match words starting and ending with specified Ukrainian vowels
	pattern := `(\s|\,|\.)[аеєиіїоуюяАЕЄИІЇОУЮЯ][а-яА-Я]*[аеєиіїоуюяАЕЄИІЇОУЮЯ](\s|\,|\.)`

	re := regexp.MustCompile(pattern)

	matches := re.FindAllString(text, -1)
	fmt.Println(matches)
}
