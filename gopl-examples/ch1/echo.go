// Eco prints its own command line arguments
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(os.Args[0] + " :")
	fmt.Println(strings.Join(os.Args[1:], " "))
	//var s, sep string

	//for _, arg := range os.Args[1:] {
	//    s += sep + arg
	//	sep = " "
	//}
	//fmt.Println(s)
}

// s+= .. :
//real    0m0.052s
//user    0m0.048s
//sys     0m0.004s

// Join:
//real    0m0.010s
//user    0m0.011s
//sys     0m0.000s
