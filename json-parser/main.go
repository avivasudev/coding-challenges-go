package main

import (
	"fmt"
	"json-parser/parser"
	"os"
)

func main() {
	// Check for exactly one argument (filename)
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	// Read the file content
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Validate JSON content
	err = parser.ValidateJSON(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
		os.Exit(1)
	}

	// Success - valid JSON
	fmt.Println("Valid JSON")
	os.Exit(0)
}