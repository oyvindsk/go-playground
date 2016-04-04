package main

import (
	"bytes"
	"fmt"
	"strings"
)

// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

// comma inserts commas in a non-negative decimal integer string.
func comma2(s string) string {
	//n := len(s)
	var buf bytes.Buffer
	var buf2 bytes.Buffer

	for i := len(s) - 1; i >= 0; i-- {
		if i%3 == 0 && i < len(s)-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte(s[i])
	}

	b := buf.Bytes()
	for i := len(b) - 1; i >= 0; i-- {
		buf2.WriteByte(b[i])
	}

	//return comma(s[:n-3]) + "," + s[n-3:]
	return buf2.String()
}

// handle floats by pretending the string ends at .
func comma3(s string) string {
	n := strings.Index(s, ".")
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

func anagram(s1, s2 string) bool {

	if len(s1) != len(s1) {
		// assume words must be the same length to be an anagram
		// not stricly correct since len counts bytes?
		return false
	}

	for _, l := range s1 {
		if !strings.Contains(s2, string(l)) {
			//fmt.Println("Not in s2:", string(l))
			return false
		}
	}

	return true
}

func anagram2(s1, s2 string) bool {

	runesSeen := make(map[rune]bool)

	for _, r := range s1 {
		runesSeen[r] = true
	}

	for _, r := range s2 {
		if !runesSeen[r] {
			fmt.Println("\tNot in s1:", string(r))
			return false
		}
		runesSeen[r] = false
	}

	// where all runes seen?
	for r, b := range runesSeen {
		if b {
			fmt.Println("\tNot in s2:", string(r))
			return false
		}
	}

	//fmt.Println(runesSeen)

	return true
}

func main() {

	fmt.Println(anagram2("abcd", "bcad"))
	fmt.Println(anagram2("foo", "bar"))
	fmt.Println(anagram2("foo", "f"))

}
