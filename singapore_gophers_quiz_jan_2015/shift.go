package main

import "fmt"

func f(a int, b uint) {
	var min = 0
	min = (a + int(b) - ((a-int(b))>>31 ^ (a - int(b))) - (a-int(b))>>31) / 2
	fmt.Printf("The min of %d and %d is %d\n", a, b, min)
}

func main() {
	f(9000, 314)
}
