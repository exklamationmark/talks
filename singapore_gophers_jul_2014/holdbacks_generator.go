package holdbacks

import (
  `strings`
)

// listen to the command channel
// when either `hb.all` or `hb.resourceID` is requested
// sends resource ids into toProcess channel
// START OMIT
func generate(command string, done <-chan bool, preloaded *plData) <-chan string {
  toProcess := make(chan string)

  go func() {
    defer close(toProcess)
    switch (command) {
    case fullCmd:
      // send a bunch of ids to channel
      sendIDs(command, toProcess, done, preloaded)
    default:
      resourceID := strings.Split(command, `.`)[1]
      select {
      case <-done:
        return
      case toProcess <- resourceID:
      }
    }
  }()

  return toProcess
}
// END OMIT

// query db for resource ids, and send the onto the channel
// the query needs to return id only
func sendIDs(resourceType string, toProcess chan<- string, done <-chan bool, preloaded *plData) {
  switch (resourceType) {
  case `full`:
    // make an array of containers & videos need to be processed, with
    // - all containers & videos explicitly mentioned in holdbacks
    //   using the res id from resToHbs
    // - for containers that was marked in resToHbs, get videos
    //   that are in normal or pending state & belong to the container
    allIDs := make([]string, 0, len(preloaded.cons) + len(preloaded.vids))
    for id, _ := range preloaded.resToHbs {
      allIDs = append(allIDs, id)
      if strings.HasSuffix(id, `c`) {
        for _, vid := range preloaded.conToVids[id] {
          allIDs = append(allIDs, vid)
        }
      }
    }

    for _, id := range allIDs {
      select {
      case toProcess <- id:
      case <-done:
        return
      }
    }
  }
}
