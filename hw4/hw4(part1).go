package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	file, err := openFile()
	if err != nil {
		fmt.Println("Error open file:", err)
		return
	}
	textLines := readData(file)
	input := inputWord(err)

	result := checkPhrase(textLines, input)
	printResult(result)
}

func readData(file *os.File) []string {
	scanner := bufio.NewScanner(file)
	var textLines []string
	for scanner.Scan() {
		line := scanner.Text()
		textLines = append(textLines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
	}
	fmt.Println("Input text")
	fmt.Println(textLines)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from file:", err)
	}
	return textLines
}

func openFile() (*os.File, error) {
	scannerFile := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the full path to the target file: ")
	scannerFile.Scan()
	filePath := scannerFile.Text()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		fmt.Println("Opened default file")
		filePath = "sky.txt"
		file, err = os.Open(filePath)
		return file, err
	}
	defer file.Close()
	fmt.Println("File opened successfully:", filePath)
	return file, nil
}

func inputWord(err error) string {
	fmt.Print("Enter word which you want find:")
	var input string
	_, err = fmt.Scanln(&input)
	fmt.Println("You entered:", input)
	return input
}

func checkPhrase(textLines []string, phrase string) []string {
	var textLineWhichRepeated []string

	for _, textLine := range textLines {
		if strings.Contains(textLine, phrase) {
			textLineWhichRepeated = append(textLineWhichRepeated, textLine)
		}
	}
	return textLineWhichRepeated
}

func printResult(textLines []string) {

	if len(textLines) == 0 {
		fmt.Printf("Word doesn't exist \n")
		return
	}
	fmt.Printf("Word exists in the lines:\n")
	for _, textLine := range textLines {
		fmt.Println(textLine)
	}

}
