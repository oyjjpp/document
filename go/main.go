package main

import (
	"log"
	"os"
)

func main() {
	os.Exit(1)
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
