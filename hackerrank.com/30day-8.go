package main

import (
	"fmt"
	"io"
)

func main() {
	var cnt int
	_, err := fmt.Scan(&cnt)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	pb := parsePhonebook(cnt)

	var name string
	//var number int

	for {
		_, err = fmt.Scan(&name)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		number, ok := pb[name]
		if ok {
			fmt.Printf("%s=%d\n", name, number)
		} else {
			fmt.Println("Not found")
		}
	}
}

func parsePhonebook(cnt int) map[string]int {

	var name string
	var number int
	pb := make(map[string]int)

	for ; cnt > 0; cnt-- {
		_, err := fmt.Scan(&name)
		if err != nil {
			fmt.Println(err)
			break
		}

		_, err = fmt.Scan(&number)
		if err != nil {
			fmt.Println(err)
			break
		}

		pb[name] = number
	}

	return pb
}
