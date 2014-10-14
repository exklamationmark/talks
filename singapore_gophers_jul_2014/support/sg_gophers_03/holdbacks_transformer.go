gackage holdbacks

import (
  customClock `github.com/101loops/clock`
  `github.com/golang/glog`
  `holdbacks_aggregator/models/captable`
  `sync`
  `holdbacks_aggregator/services/queue`
  `os`
  `strings`
)

type hbData struct {
  sync.RWMutex
  tb *captable.Captable
}

const (
  bulkSliceSize = 10
)

var (
  clock customClock.Clock
  pre *plData
  capTable *hbData
)


func init() {
  clock = customClock.New()
  pre = preload()
  capTable = &hbData {
    tb: captable.New(),
  }
}

// START OMIT
func Process(command string, maxGoRoutine int, done <-chan bool) <-chan *queue.Message {
  // pre = preloaded data, capTable == holdbacks data
  toProcess := generate(command, done, pre)

  workers := make([]<-chan string, maxGoRoutine)
  for i := 0; i < maxGoRoutine; i++ {
    workers[i] = worker(toProcess, capTable, done, pre)
  }

  // monitor progress
  doneSoFar := 0
  for res := range merge(done, workers) {
    doneSoFar += 1
    glog.Info(`processed: (`, doneSoFar, `) `, res )
  }

  // additional, global steps
  // some more processing

  queueMsgs := generateQueueMsg(capTable, done, bulkSliceSize)
  return queueMsgs
}
// END OMIT
