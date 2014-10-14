package main

import (
  `time`
  `math/rand`
  `fmt`
)

// START OMIT
func main() {
  first := fetch(`Gopher's guide to universe`)
  second := fetch(`Gopher 101`)
  book1 := <- first
  book2 := <- second
  fmt.Println(`burn:`, book1, `and`, book2)
}

func fetch(name string) <-chan string {
  out := make(chan string)
  go func() {
    defer close(out)
    time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
    out <- name
  }()
  return out
}
// END OMIT
