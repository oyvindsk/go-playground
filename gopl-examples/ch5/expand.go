package main

import (
	"fmt"
	"strings"
)

func main() {

	s := "hello $foo, hello thisis$allowed $upperacseme $please"

	//upcase := func(s string) string {
	//	return strings.ToUpper(s)
	//}

	fmt.Println(expand(s, strings.ToUpper))

}

func expand(str string, f func(string) string) string {
	// split it
	s := strings.Fields(str)
	for i, w := range s {
		// check for $
		if strings.HasPrefix(w, "$") {
			// run f()
			s[i] = f(w[1:]) // ignore $
		}
	}

	return strings.Join(s, " ")

}
