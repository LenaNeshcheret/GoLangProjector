package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.OpenFile("/Users/bigmag/GoLangProjector/hw12/KeyPassword.txt", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("The file to save doesn't answer", err)
	}
	defer file.Close()

	save(file)
	printPasswordKey(getAllPasswordKey(file))
	getPasswordByKey(file)
}

func save(file *os.File) {

	fmt.Println("Enter key for password:")
	var PasswordKey string
	_, err := fmt.Scan(&PasswordKey)
	if err != nil {
		fmt.Println("Error reading key for password:", err)
	}
	keys := getAllPasswordKey(file)
	for _, key := range keys {
		if key == PasswordKey {
			fmt.Printf("Please use another key for password, %s is already used\n", PasswordKey)
			return
		}
	}

	fmt.Println("Enter password:")
	var Password string
	_, err = fmt.Scan(&Password)
	if err != nil {
		fmt.Println("Error reading password:", err)
	}

	_, err = file.WriteString(PasswordKey + " " + Password + "\n")
	if err != nil {
		fmt.Println("Error write to file:", err)
	}
}

func getAllPasswordKey(file *os.File) []string {
	var keys []string
	resetReadToStart(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		keys = append(keys, parts[0])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return keys
}
func printPasswordKey(pats []string) {
	fmt.Println("All PasswordKeys:")
	for _, p := range pats {
		fmt.Println(p)
	}
}

func getPasswordByKey(file *os.File) {
	fmt.Println("Enter password keys:")
	var PasswordKey string
	_, err := fmt.Scan(&PasswordKey)
	if err != nil {
		fmt.Println("Error reading key for password:", err)
	}
	resetReadToStart(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		for _, passKey := range parts {
			if PasswordKey == passKey {
				fmt.Println("This is your password:")
				fmt.Println(parts[1])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

}

func resetReadToStart(file *os.File) {
	_, err := file.Seek(0, 0)
	if err != nil {
		fmt.Println("The file to save doesn't answer", err)
	}
}
