package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func countBytes(file *os.File) int {

	reader := bufio.NewReader(file)

	count := 0

	buffer := make([]byte, 1024) // Read in 1KB chunks

	for {
		chunk_count, err := reader.Read(buffer)
		count += chunk_count

		if err != nil {
			if err == io.EOF {
				break // End of file
			}
			log.Fatal(err)
		}
	}

	return count
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
		log.SetOutput(os.Stderr)
		log.Printf("reading input to count lines: %v\n", err)
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
		log.SetOutput(os.Stderr)
		log.Printf("reading input to count words: %v\n", err)
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
		log.SetOutput(os.Stderr)
		log.Printf("reading input to count chars: %v\n", err)
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

	if flag.NFlag() == 1 {

		if *c {
			count := countBytes(file)
			fmt.Printf("%d %s\n", count, filePath)

		} else if *l {
			count := countLines(file)
			fmt.Printf("%d %s\n", count, filePath)

		} else if *w {
			count := countWords(file)
			fmt.Printf("%d %s\n", count, filePath)

		} else if *m {
			count := countChars(file)
			fmt.Printf("%d %s\n", count, filePath)

		}
	} else if flag.NFlag() == 0 {

		byte_count := countBytes(file)
		file.Seek(0, io.SeekStart)
		line_count := countLines(file)
		file.Seek(0, io.SeekStart)
		word_count := countWords(file)

		fmt.Printf("%d %d %d %s\n", byte_count, line_count, word_count, filePath)

	}
}
