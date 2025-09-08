# Word Count Tool (wc) in Go

## Overview

This command-line tool is a Go-based implementation of the basic functionality of the Unix `wc` (word count) command. It is designed to count the number of lines, words, bytes, and characters from a single file or from standard input.

## How to Run

You can run the script using the Go toolchain.

### Running Directly

To run the script without building an executable, use `go run`:
```bash
go run wc/wc.go test.txt
```

### Building the Executable

For repeated use, you can build a compiled executable:
```bash
go build wc/wc.go
```
This will create an executable file named `wc` in the current directory. You can then run it like this:
```bash
./wc test.txt
```

### Reading from Standard Input

The tool can also process data piped from other commands if no file name is provided:
```bash
cat test.txt | go run wc/wc.go
```

### Functionality & Flags

The script's behavior changes based on the flags provided.

*   **Default (No Flags)**: If you don't specify any flags, the tool will print the line, word, and byte counts for the given file.
    ```bash
    go run wc/wc.go test.txt
    ```

*   **Flags**: You can specify one of the following flags to get a specific count:
    *   `-c`: Prints the number of bytes in the file.
    *   `-l`: Prints the number of lines in the file.
    *   `-w`: Prints the number of words in the file.
    *   `-m`: Prints the number of characters in the file.

    **Example:** To get only the character count:
    ```bash
    go run wc/wc.go -m test.txt
    ```