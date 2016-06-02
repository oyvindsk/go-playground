package main

import (
	"bytes"
	"fmt"
)

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
}

func main() {

	var x, y IntSet
	y.AddAll(1, 2, 3, 144, 9)
	x.AddAll(2, 144, 8)
	fmt.Println(x.String())
	fmt.Println(y.String())
	//x.IntersectWith(&y)
	//fmt.Println(x.String())
	for _, w := range y.Elems() {
		fmt.Println("\t", w)
	}
	//x.UnionWith(&y)
	//fmt.Println(x.String())           // "{1 9 42 144}"
	//fmt.Println(x.Has(9), x.Has(123)) // "true false"
}

func (s *IntSet) Elems() []int {

	var elems []int
	for i, word := range s.words {
		if word == 0 {
			continue
		}

		for j := 0; j < 64; j++ {
			if (word & (1 << uint(j))) > 0 {
				fmt.Println("FOUND:", i, j, i*64+j)
				elems = append(elems, i*64+j)
			}
		}
	}
	return elems
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
			//} else {
			//	s.words = append(s.words, tword)
		}
	}

	// cut the extra words in our words
	if len(s.words) > len(t.words) {
		s.words = s.words[:len(t.words)]
	}

}

// AddAll adds a bunch of non-negative values to the set.
func (s *IntSet) AddAll(xx ...int) {
	for _, x := range xx {
		s.Add(x)
	}
}

func (s *IntSet) Len() int {

	var cnt int

	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				cnt++
			}
		}
	}

	return cnt
}

// Remove removed the non-negative value x from the set.
// TODO: Handle non-existing values, now they're just added
func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] ^= 1 << bit // Xor the bit out
}

func (s *IntSet) Clear() {
	for i := range s.words {
		s.words[i] = 0
		//word &= 1 << 63
	}
}

func (s *IntSet) Copy() *IntSet {
	is := new(IntSet) // returns a pointer to  a new IntSet
	for _, w := range s.words {
		is.words = append(is.words, w)
	}
	return is

}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}
