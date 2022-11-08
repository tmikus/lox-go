package main

import (
	"bufio"
	"fmt"
	"os"
)

var hasError bool = false

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage jlox [script]")
		os.Exit(64)
	}
	if len(args) == 1 {
		exitCode := runFile(args[0])
		os.Exit(exitCode)
	} else {
		runPrompt()
	}
}

//func printError(line int, message string) {
//	reportError(line, "", message)
//}
//
//func reportError(line int, where string, message string) {
//	fmt.Println("[line", line, "] Error", where, ":", message)
//}

func run(source string) int {
	scanner := NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
	return 0
}

func runFile(filePath string) int {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file", err)
		return 1
	}
	fileContent := string(fileBytes)
	fmt.Println("Running file: ", fileContent)
	return 0
}

func runPrompt() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return 1
		}
		// Strip newline at the end
		text = text[:len(text)-1]
		run(text)
	}
}
