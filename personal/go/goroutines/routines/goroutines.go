package main

import (
	"fmt"
	"math/rand"
	"time"
)

func hello(msg string) {
	fmt.Println(msg + "- goroutine")
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Second)
}

func main() {
	go hello("hello 1")
	go hello("hello 2")

	time.Sleep(1 * time.Second)
	fmt.Println("Chamada normal")
}
