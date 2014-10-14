package main

import (
  `flag`
  `holdbacks_aggregator/transformers/holdbacks`
  `holdbacks_aggregator/helpers/config`
  `holdbacks_aggregator/services/queue`
  `github.com/golang/glog`
  `os`
)

var (
  aggregatorJob, aggregatorResource string
)

func init() {
  flag.StringVar(&aggregatorJob, `j`, `full`, `type of job (full/partial)`)
  flag.StringVar(&aggregatorResource, `r`, `holdbacks`, `resources to process (holdbacks)`)
}

func main() {
  flag.Parse()

  var output <-chan *queue.Message
  done := make(chan bool)
  defer func() {
    glog.Info(`closing done`)
    close(done)
  }()

  switch(aggregatorResource) {
  case `holdbacks`:
    output = holdbacks.Process(aggregatorJob, config.Config.Routines, done)
  }

  for msg := range output {
    switch(msg.Action) {
    case `update`:
      queue.Client.UpdateMessage(msg.Resource, msg.ID, msg.Data.(string), `services.hyperion2`)
      glog.Info(`enqueue: `, msg.Resource, ` `, msg.ID, ` `, msg.Data.(string), ` services.hyperion2`)
    default:
      glog.Error(`can't send this message `, msg)
    }

  }
  f.Sync()
}
