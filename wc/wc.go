package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func countBytes(file *os.File) int64 {

	fileInfo, err := file.Stat()

	if err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return fileInfo.Size()
}

func countLines(file *os.File) int {

	count := 0

	// Read the file contents
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return count

}

func countWords(file *os.File) int {

	count := 0

	// Read the file contents
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return count

}

func countChars(file *os.File) int {

	count := 0

	// Read the file contents
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
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

	var file *os.File
	var file_err error
	var filePath string

	if len(args) == 0 {

		file = os.Stdin

	} else if len(args) == 1 {
		filePath = args[0]

		// Open the file
		file, file_err = os.Open(filePath)

		if file_err != nil {
			log.Fatalf("Failed to open the file: %v", file_err)
		}
		defer file.Close()
	}

	// fmt.Println(filePath)

	if flag.NFlag() == 1 {

		if *c {
			count := countBytes(file)
			fmt.Printf("%d %s", count, filePath)

		} else if *l {
			count := countLines(file)
			fmt.Printf("%d %s", count, filePath)

		} else if *w {
			count := countWords(file)
			fmt.Printf("%d %s", count, filePath)

		} else if *m {
			count := countChars(file)
			fmt.Printf("%d %s", count, filePath)

		}
	} else if flag.NFlag() == 0 {

		byte_count := countBytes(file)
		line_count := countLines(file)
		word_count := countWords(file)

		fmt.Printf("%d %d %d %s", byte_count, line_count, word_count, filePath)

	}
}
