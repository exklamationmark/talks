package main

import (
	"database/sql"
	"log"
	"time"
)

func sustained(db *sql.DB) {
	var limit int
	if *maxOpen > 0 {
		limit = int(float64(*maxOpen) * (*multiplier))
	} else {
		log.Printf("[%d] No connection limit. Set base concurrency= 100 goroutines\n", time.Now().UTC().UnixNano())
		limit = *ceil
	}

	end := make(chan struct{})
	for i := 0; i < limit; i++ {
		id := i
		go sustainedWorker(id, db, *interval, end)
	}

	<-time.After(*duration)
	close(end)
	time.Sleep(time.Second)
}

func sustainedWorker(id int, db *sql.DB, interval time.Duration, end <-chan struct{}) {
	for {
		select {
		case <-end:
			return
		default:
			doWork(id, db)
			time.Sleep(interval)
		}
	}
}
