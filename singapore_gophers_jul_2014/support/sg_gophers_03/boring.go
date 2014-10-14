package main

import (
 `fmt`
 `time`
 `math/rand`
)

func main() {
  boring(`test`)
}

// START OMIT
func boring(msg string) {
  for i := 0; ; i++ {
    fmt.Println(msg, i)
    time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
  }
}
// END OMIT
