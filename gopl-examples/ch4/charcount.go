// Charcount computes counts of Unicode characters.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int)     // counts of Unicode characters
	catCount := make(map[string]int) // count the category each rune belongs to (letter, number, printable etc)
	var utflen [utf8.UTFMax + 1]int  // count of lengths of UTF-8 encodings
	invalid := 0                     // count of invalid UTF-8 characters
	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++

		if unicode.IsControl(r) {
			catCount["control"]++
		}
		if unicode.IsDigit(r) {
			catCount["digit"]++
		}
		if unicode.IsGraphic(r) {
			catCount["graphic"]++
		}
		if unicode.IsLetter(r) {
			catCount["letter"]++
		}
		if unicode.IsLower(r) {
			catCount["lower case"]++
		}
		if unicode.IsMark(r) {
			catCount["mark"]++
		}
		if unicode.IsNumber(r) {
			catCount["number"]++
		}
		if unicode.IsPrint(r) {
			catCount["printable"]++
		}
		if !unicode.IsPrint(r) {
			catCount["non-printable"]++
		}
		if unicode.IsPunct(r) {
			catCount["punct"]++
		}
		if unicode.IsSpace(r) {
			catCount["space"]++
		}
		if unicode.IsSymbol(r) {
			catCount["symbol"]++
		}
		if unicode.IsTitle(r) {
			catCount["title case"]++
		}
		if unicode.IsUpper(r) {
			catCount["upper case"]++
		}

	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}

	fmt.Printf("category\tcount\n")
    for k,v := range catCount {
        fmt.Printf("%q\t%d\n", k, v)
    }

	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}
