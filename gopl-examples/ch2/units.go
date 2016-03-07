package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/oyvindsk/go-playground/gopl-examples/ch2/tempconv"
)

type Meters float64
type Feet float64

const (
	MetersInAFoot = 0.3048
)

func main() {
	for _, arg := range os.Args[1:] {
		u, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "units: %v\n", err)
			os.Exit(1)
		}
		f := tempconv.Fahrenheit(u)
		c := tempconv.Celsius(u)
		fmt.Printf("%s = %s, %s = %s\n",
			f, tempconv.FToC(f), c, tempconv.CToF(c))

		m := Meters(u)
		fe := Feet(u)
		fmt.Printf("%s = %s, %s = %s\n", fe, FToM(fe), m, MToF(m))
	}

}

func MToF(m Meters) Feet { return Feet(m / MetersInAFoot) }

func FToM(f Feet) Meters { return Meters(f * MetersInAFoot) }

func (m Meters) String() string { return fmt.Sprintf("%g Meters", m) }
func (f Feet) String() string {
	switch {
	case f == 1:
		return fmt.Sprintf("%g Foot", f)
	default:
		return fmt.Sprintf("%g Feet", f)
	}
}
