package main

import "fmt"

func f(a, b int) {
	var min = (map[bool]int{true: b, false: a})[a > b]
	fmt.Printf("The min of %d and %d is %d\n", a, b, min)
	fmt.Printf("The min of %d and %d isn't %d\n", a, b, min)
}

func main() {
	f(9000, 314)
}
