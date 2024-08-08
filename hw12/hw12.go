package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		printUsage()
		return
	}

	passwordManager := NewPasswordManager(os.Args[1])
	if err := passwordManager.OpenFile(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer passwordManager.CloseFile()

	if len(os.Args) < 3 {
		printUsage()
		return
	}

	command := os.Args[2]
	passwordManager.ExecuteCommand(command, os.Args[3:])
}

type passwordManager struct {
	filePath string
	file     *os.File
}

func NewPasswordManager(filePath string) *passwordManager {
	return &passwordManager{filePath: filePath}
}

func (pM *passwordManager) OpenFile() error {
	file, err := os.OpenFile(pM.filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	pM.file = file
	return nil
}

func (pM *passwordManager) CloseFile() {
	if pM.file != nil {
		pM.file.Close()
	}
}

func (pM *passwordManager) ExecuteCommand(command string, args []string) {
	switch command {
	case "getAll":
		if len(args) != 0 {
			fmt.Println("Usage: getAll")
			return
		}
		pM.getAllPasswordKeys()
	case "get":
		if len(args) != 1 {
			fmt.Println("Usage: get <key>")
			return
		}
		pM.getPasswordByKey(args[0])
	case "save":
		if len(args) != 2 {
			fmt.Println("Usage: save <key> <password>")
			return
		}
		pM.savePassword(args[0], args[1])
	default:
		printHelp(command)
	}
}

func (pM *passwordManager) getAllPasswordKeys() {
	keys := pM.readAllPasswordKeys()
	printPasswordKeys(keys)
}

func (pM *passwordManager) getPasswordByKey(key string) {
	resetReadToStart(pM.file)
	scanner := bufio.NewScanner(pM.file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if parts[0] == key {
			fmt.Println("Password for key:", key)
			fmt.Println(parts[1])
			return
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	fmt.Println("Key not found:", key)
}

func (pM *passwordManager) savePassword(key, password string) {
	keys := pM.readAllPasswordKeys()
	for _, k := range keys {
		if k == key {
			fmt.Printf("Key '%s' is already used. Please use another key.\n", key)
			return
		}
	}
	_, err := pM.file.WriteString(key + " " + password + "\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func (pM *passwordManager) readAllPasswordKeys() []string {
	var keys []string
	resetReadToStart(pM.file)
	scanner := bufio.NewScanner(pM.file)
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

func resetReadToStart(file *os.File) {
	_, err := file.Seek(0, 0)
	if err != nil {
		fmt.Println("Error resetting file read position:", err)
	}
}

func printPasswordKeys(keys []string) {
	fmt.Println("All Password Keys:")
	for _, key := range keys {
		fmt.Println(key)
	}
}

func printHelp(command string) {
	fmt.Println("Unknown command:", command)
	printUsage()
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" <file_path> getAll")
	fmt.Println(" <file_path> get <key>")
	fmt.Println(" <file_path> save <key> <password>")
}
