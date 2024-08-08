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

	mapByWords := initTextByWords(textLines)
	input := inputWord(err)
	printAllLinesByWorld(input, mapByWords)
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
		filePath = "testHw5.txt"
		file, err = os.Open(filePath)
		return file, err
	}
	defer file.Close()
	fmt.Println("File opened successfully:", filePath)
	return file, nil
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
	return textLines
}

func inputWord(err error) string {
	fmt.Print("Enter word which you want find:")
	var input string
	_, err = fmt.Scanln(&input)
	fmt.Println("You entered:", input)
	return input
}

func initTextByWords(textLines []string) map[string][]string {

	initMap := make(map[string][]string)

	for _, textLine := range textLines {
		splitWords := strings.Fields(textLine)
		for _, splitWord := range splitWords {

			initMap[splitWord] = append(initMap[splitWord], textLine)
		}
	}
	return initMap
}

func printAllLinesByWorld(word string, textLines map[string][]string) {
	vals := textLines[word]
	if len(vals) == 0 {
		fmt.Printf("Word [%s] doesn't exist \n", word)
		return
	}
	fmt.Printf("Word [%s] exists in the lines:\n", word)
	for _, val := range vals {
		fmt.Println(val)
	}
}
