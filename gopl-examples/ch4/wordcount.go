package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// An artificial input source.
	scanner := bufio.NewScanner(os.Stdin)
	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)
	// Count the words.
	count := 0
	words := make(map[string]int)

	for scanner.Scan() {
		count++
		words[scanner.Text()]++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	fmt.Printf("Total: %d\n", count)

	pairList := rankByWordCount(words)
	for i, p := range pairList {
		fmt.Printf("%s\t%d\n", p.Key, p.Value)
		if i > 8 {
			break
		}
	}
}

// Sort the Pair type by value
func rankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, 0, len(wordFrequencies))
	for k, v := range wordFrequencies {
		pl = append(pl, Pair{k, v})
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// Each key, val - Easy to compare
type Pair struct {
	Key   string
	Value int
}

// List of Pair's  - Teh type we want to sort with custom funcs
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
