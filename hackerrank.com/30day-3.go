package main

import "fmt"

func main() {
	var cost float64
	var tip, tax int

	cnt, err := fmt.Scan(&cost)
	if err != nil {
		panic(fmt.Sprintf("Only read: %d, err: %v", cnt, err))
	}

	cnt, err = fmt.Scan(&tip)
	if err != nil {
		panic(fmt.Sprintf("Only read: %d, err: %v", cnt, err))
	}

	cnt, err = fmt.Scan(&tax)
	if err != nil {
		panic(fmt.Sprintf("Only read: %d, err: %v", cnt, err))
	}

	//fmt.Println(cost, tip, tax)

	total := cost + (cost * (float64(tip) / 100.0)) + (cost * (float64(tax) / 100.0))
	fmt.Printf("The total meal cost is %.0f dollars.\n", total)

}
