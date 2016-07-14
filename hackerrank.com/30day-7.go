package main

import "fmt"

func main() {
	var cnt int
	_, err := fmt.Scan(&cnt)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var i int
	arr := make([]int, cnt)

	for ; cnt > 0; cnt-- {
		_, err = fmt.Scan(&i)
		if err != nil {
			fmt.Println(err)
			break
		}
		arr[cnt-1] = i

	}

	for _, j := range arr {
		fmt.Printf("%d ", j)
	}
	fmt.Println()
}
