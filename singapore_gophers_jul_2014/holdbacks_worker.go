// worker takes in a resource from the channel and run baseHoldback for it
package holdbacks

import (
  `sync`
  `strings`
)

func worker(input <-chan string, d *hbData, done <-chan bool, preloaded *plData) <-chan string {
  output := make(chan string)

  go func() {
    defer close(output)
    for res := range input {
      isVideo := strings.HasSuffix(res, `v`)
      holdbacks := hbsFor(res, preloaded, isVideo)
      baseHoldbacks(res, holdbacks, preloaded, isVideo, d)
      select {
      case output <- res:
      case <-done:
      }
    }
  }()

  return output
}

// START OMIT
func merge(done <-chan bool, processed []<-chan string) <-chan string {
    var wg sync.WaitGroup
    out := make(chan string)

    addToOutput := func(single <-chan string) {
      defer wg.Done()
      for msg := range single {
        select {
        case out <- msg:
        case <-done:
          return
        }
      }
    }

    wg.Add(len(processed))
    for _, c := range processed {
      go addToOutput(c)
    }

    go func() {
      wg.Wait()
      close(out)
    }()
    return out
}
// END OMIT
