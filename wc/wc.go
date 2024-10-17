package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func countBytes(filePath string) int64 {
	// Get file information
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("The file %s does not exist", filePath)
		} else {
			log.Fatalf("Failed to get file information: %v", err)
		}
	}

	return fileInfo.Size()
}

func countLines(filePath string) int {
	// Open the file
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("Failed to open the file: %v", err)
	}
	defer file.Close()

	count := 0

	// Read the file contents
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		count++
	}

	return count

}

func countWords(filePath string) int {
	// Open the file
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("Failed to open the file: %v", err)
	}
	defer file.Close()

	count := 0

	// Read the file contents
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		count++
	}

	return count

}

func countChars(filePath string) int {
	// Open the file
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("Failed to open the file: %v", err)
	}
	defer file.Close()

	count := 0

	// Read the file contents
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		count++
	}

	return count

}

func main() {

	// Define flags
	c := flag.Bool("c", false, "print no of bytes in file")
	l := flag.Bool("l", false, "print no of lines in file")
	w := flag.Bool("w", false, "print no of words in file")
	m := flag.Bool("m", false, "print no of characters in file")

	// Parse flags
	flag.Parse()

	// The remaining arguments after flags are parsed
	args := flag.Args()

	fmt.Printf("%d %d\n", len(args), flag.NFlag())

	var filePath string

	if len(args) == 0 {

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			filePath = scanner.Text()
		}

		// Check for errors
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

		fmt.Println(filePath)

	}

	fmt.Println(filePath)

	if flag.NFlag() == 1 {
		if filePath == "" {
			filePath = args[0]
		}

		if *c {
			count := countBytes(filePath)
			fmt.Printf("%d %s", count, filePath)

		} else if *l {
			count := countLines(filePath)
			fmt.Printf("%d %s", count, filePath)

		} else if *w {
			count := countWords(filePath)
			fmt.Printf("%d %s", count, filePath)

		} else if *m {
			count := countChars(filePath)
			fmt.Printf("%d %s", count, filePath)

		}
	} else if flag.NFlag() == 0 {
		if filePath == "" {
			filePath = args[0]
		}
		byte_count := countBytes(filePath)
		line_count := countLines(filePath)
		word_count := countWords(filePath)

		fmt.Printf("%d %d %d %s", byte_count, line_count, word_count, filePath)

	}
}
