package main

import (
	"fmt"

	"github.com/oyvindsk/go-playground/gopl-examples/ch2/tempconv"
)

func main() {
	fmt.Println(tempconv.KtoC(1000))
	fmt.Println(tempconv.CtoK(-0))
}
