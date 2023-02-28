package main

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {

	testcases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}

	for _, tc := range testcases {
		rev := Reverse(tc.in)

		if rev != tc.want {
			t.Errorf("Reverse: got %q, want: %q", rev, tc.want)
		}
	}
}

func FuzzReverse(f *testing.F) {

	testcases := []string{"Hello, World", " ", "!12345"}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)
		revrev := Reverse(rev)
		if orig != revrev {
			t.Errorf("Before: %q, after: %q", orig, revrev)
		}

		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid utf8 string: %q", rev)
		}
	})
}
