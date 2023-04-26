package main

import "fmt"

func main() {
	var name interface{} = "frank"
	// a, ok := name.(int)
	// fmt.Println(a, ok)
	a := name.(int)
	fmt.Println(a)

	// 5 / 4
}
