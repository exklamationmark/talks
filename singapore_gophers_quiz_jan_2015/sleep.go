package main

import (
	"fmt"
	"time"
)
func f(a, b int) {
	var min = 0
	c := make(chan int)
	go func() {
		time.Sleep(time.Duration(a))
		c <- a
	}()
	go func(){
		time.Sleep(time.Duration(b))
		c <- b
	}()
	min = <-c
	fmt.Printf("The min of %d and %d is %d\n", a, b, min)
}

func main() {
	f(9000, 314)
}
