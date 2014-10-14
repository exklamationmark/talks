package main

import (
 `fmt`
 `time`
 `math/rand`
)

var (
  names = map[int]string{0: `alpha`, 1: `beta`, 2: `gamma`, 3: `delta`}
)

// START OMIT
func main() {
  for i := 0; i < 4; i++ {
    go boring(names[i])
  }
  time.Sleep(time.Second * 30)
}
// END OMIT

func boring(msg string) {
  for i := 0; ; i++ {
    fmt.Println(msg, i)
    time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
  }
}
