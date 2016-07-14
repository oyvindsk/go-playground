package main

import "fmt"

func main() {
	var strCnt int
	_, err := fmt.Scan(&strCnt)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var s string

	for ; strCnt > 0; strCnt-- {
		_, err = fmt.Scan(&s)
		if err != nil {
			fmt.Println(err)
			break
		}

		// do two passes over the string so we do not have to store anything
		for i, r := range s {
			if i%2 == 0 {
				fmt.Print(string(r))
			}
		}

		fmt.Print(" ")

		for i, r := range s {
			if i%2 == 1 {
				fmt.Print(string(r))
			}
		}
		fmt.Print("\n")
	}
}
