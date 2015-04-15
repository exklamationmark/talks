package main

import (
	`encoding/json`
	`log`
	`net/http`
	_ `net/http/pprof`
	`time`
)

type Gopher struct {
	BirthDay int64 `json:"birth_date"`
}

var (
	errNoGopher = []byte(`{"failed":"no gopher for you"}`)
)

func main() {
	http.HandleFunc(`/`, handler)
	log.Fatal(http.ListenAndServe(`:9090`, nil))
}

func handler(resp http.ResponseWriter, req *http.Request) {
	body, err := json.Marshal(NewGopher())
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write(errNoGopher)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write(body)
}

func NewGopher() *Gopher {
	gopher := Gopher{time.Now().Unix()}
  return &gopher
}
