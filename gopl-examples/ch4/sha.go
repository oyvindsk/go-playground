package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {

	shaA := sha256.Sum256([]byte("x"))
	shaB := sha256.Sum256([]byte("X"))
	diff := compareSha(&shaA, &shaB)
	fmt.Println("Diff between A and B:", diff)

}

func compareSha(shaA *[32]byte, shaB *[32]byte) int {
	var diff int

	for i, _ := range shaA {
		bA := shaA[i]
		bB := shaB[i]

		for j := uint(0); j < 8; j++ {
			// Loop each bit, start with the last (rightmost) one
			//fmt.Printf("A:%8b B:%8b\n", ((bA>>j)&1), ((bB>>j)&1))
			if ((bA >> j) & 1) != ((bB >> j) & 1) {
				diff++
			}
		}

	}

	return diff
}
