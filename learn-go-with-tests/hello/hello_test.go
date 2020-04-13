package main

import "testing"

func TestHello(t *testing.T) {

	assertCorrectMsg := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got: %q, want: %q", got, want)
		}
	}

	t.Run("saying hello to ppl", func(t *testing.T) {
		got := Hello("Chris", "")
		want := "Hello, Chris"
		assertCorrectMsg(t, got, want)

	})

	t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertCorrectMsg(t, got, want)
	})

	t.Run("in spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMsg(t, got, want)
	})

	t.Run("in french", func(t *testing.T) {
		got := Hello("Frank", "French")
		want := "Bonjour, Frank"
		assertCorrectMsg(t, got, want)
	})
}
