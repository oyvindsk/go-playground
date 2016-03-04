package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	filenames := make(map[string]map[string]bool) // stores string double, ufgh.. maybe use a struct in counts instead of an int?
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, filenames)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, filenames)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
			fmt.Println("in files:")
			for v := range filenames[line] {
				fmt.Printf("\t%s\n", v)
			}
		}
	}
}
func countLines(f *os.File, counts map[string]int, filenames map[string]map[string]bool) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		counts[line]++
		if filenames[line] == nil {
			filenames[line] = make(map[string]bool)
		}
		filenames[line][f.Name()] = true
	}
	// NOTE: ignoring potential errors from input.Err()
}
