package main

import "fmt"

func main() {
	ch := make(chan struct{})
	ch <- struct{}{}

	data := <-ch
	fmt.Println(data)
}

// 7164

// 4500+2500+8229.07
