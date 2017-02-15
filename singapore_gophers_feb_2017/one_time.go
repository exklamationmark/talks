package main

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

func oneTime(db *sql.DB) {
	wg := &sync.WaitGroup{}

	var limit int
	if *maxOpen > 0 {
		limit = int(float64(*maxOpen) * (*multiplier))
	} else {
		log.Printf("[%d] No connection limit. Send constant load\n", time.Now().UTC().UnixNano())
		limit = 2000
	}

	for i := 0; i < limit; i++ {
		wg.Add(1)
		id := i
		go oneTimeWorker(id, wg, db)
	}

	wg.Wait()
}

func oneTimeWorker(id int, wg *sync.WaitGroup, db *sql.DB) {
	doWork(id, db)
	wg.Done()
}
