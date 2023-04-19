package main

import "log"

func main() {
	b := new(BWM)
	b.start()
	b.use()
}

type Car struct{}

func (c *Car) use() {
	log.Println("car useA")
}

func (c Car) start() {
	log.Println("start")
}

type BWM struct {
	Car
}

func (b BWM) use() {
	log.Println("bmw use")
}
