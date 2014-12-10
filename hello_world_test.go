package main

import (
	"testing"
)

func TestGravity(t *testing.T) {
	two := 1 + 1
	if two != 2 {
		t.Error("Excpected 2, got: ", two)
	}
}
