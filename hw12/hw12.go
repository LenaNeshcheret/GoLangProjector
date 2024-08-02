package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Not enough arguments")
		printUsage()
		return
	}
	fromFile := os.Args[1]
	file, err := os.OpenFile(fromFile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("The file to save doesn't answer", err)
	}
	defer file.Close()

	command := os.Args[2]

	switch command {
	case "getAll":
		if len(os.Args) != 3 {
			fmt.Println("Usage: getAll")
			return
		}
		getAllPasswordKey(file)
	case "get":
		if len(os.Args) != 4 {
			fmt.Println("Usage: get <key>")
			return
		}
		key := os.Args[3]
		getPasswordByKey(file, key)
	case "save":
		if len(os.Args) != 5 {
			fmt.Println("Usage: save <key> <password>")
			return
		}
		key := os.Args[3]
		password := os.Args[4]
		save(file, key, password)
	default:
		printHelp(command)
	}
}

func printHelp(command string) {
	fmt.Println("Unknown command:", command)
	printUsage()
}
func printUsage() {
	fmt.Println("Write the path to the file as the first parameter. After that, use one of the functions:")
	fmt.Println(" getAll")
	fmt.Println(" get <key>")
	fmt.Println(" save <key> <password>")
}

func save(file *os.File, key string, password string) {

	keys := getAllPasswordKey(file)
	for _, k := range keys {
		if k == key {
			fmt.Printf("Please use another key for password, %s is already used\n", key)
			return
		}
	}
	_, err := file.WriteString(key + " " + password + "\n")
	if err != nil {
		fmt.Println("Error write to file:", err)
		return
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
	printPasswordKey(keys)
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

func getPasswordByKey(file *os.File, key string) {
	resetReadToStart(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		for _, passKey := range parts {
			if key == passKey {
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
