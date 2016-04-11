package main

import (
	"bytes"
	"fmt"
	"unicode"
)

func main() {

	//a := [...]int{0, 1, 2, 3, 4, 5}
	//reversea(&a)
	//rotate(a[:], 4)

	//a := [...]string{"foo", "bar", "bar", "baz"}
	//s := duprem(a[:])
	//fmt.Println(a)
	//fmt.Println(s)

	//b := []byte{'a', ' ', 'b', ' ', '\n', ' ', 'c'}
	//fmt.Printf("%s\n%v\n", b, b)
	////  0  1   2  3  4  5  6
	//// [97 32 98 32 10 32 99]
	//b = spacerem(b)
	//fmt.Printf("%s\n%v\n", b, b)

	//b := []byte{'a', 'b', 'c', 'æ', 'ø', 'å'}
	b := []byte(string([]rune{'a', 'b', 'c', 'æ', 'ø', 'å'}))
	fmt.Printf("%s\n%v\n", string(b), b)
	b = reverseB(b)
	fmt.Printf("%s\n%v\n", b, b)

}

// reverse reverses a slice of runes
func reverseB(b []byte) []byte {
	r := bytes.Runes(b)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        fmt.Println("Moving:", r[i], r[j])
		r[i], r[j] = r[j], r[i]
	}
	fmt.Printf("%s\n%v\n", string(r), r)
    b = []byte(string(r))
    return b
}

func spacerem(b []byte) []byte {
	for i := range b {
		fmt.Println("i:", i)
		if i < len(b)-1 && unicode.IsSpace(rune(b[i])) {
			// at least one unicode space, make it ascii
			b[i] = ' '
			fmt.Println("First Space at:", i)

			// any following unicode spaces we should squash?
			j := i + 1
			for j < len(b) && unicode.IsSpace(rune(b[j])) {
				fmt.Println("Also Space at:", j)
				j++
			}
			fmt.Println("i:", i, "j:", j)
			copy(b[i+1:], b[j:])
			b = b[:len(b)-(j-i)+1]
		}
	}
	return b

}

func duprem(s []string) []string {
	// broken if i==i+1==i+2 ? Se spacerem
	for i := range s {
		if i < len(s)-1 && s[i] == s[i+1] {
			fmt.Println("Dup:", s[i])
			copy(s[i:], s[i+1:])
			s = s[:len(s)-1]
		}
	}

	return s
}

func rotate(s []int, n int) []int {
	tmp := make([]int, len(s))
	copy(tmp, s)
	for i, _ := range s {
		from := (i + n) % len(s)
		fmt.Printf("From %d  to  %d\n", from, i)
		s[i] = tmp[from]
	}

	fmt.Println(tmp, s)
	return s

}

func reversea(a *[6]int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// reverse reverses a slice of ints in place.
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
