package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	wrkSustained = "sustained"
	wrkOneTime   = "one_time"
)

var (
	worker map[string]func(*sql.DB) = map[string]func(*sql.DB){
		wrkOneTime:   oneTime,
		wrkSustained: sustained,
	}
)

var (
	workload   *string
	maxOpen    *int
	maxIdle    *int
	interval   *time.Duration
	duration   *time.Duration
	multiplier *float64
)

func main() {
	workload = flag.String("workload", wrkOneTime, "describe the workload")
	maxOpen = flag.Int("max-open", 0, "max open db conn")
	maxIdle = flag.Int("max-idle", 0, "max idle db conn")
	interval = flag.Duration("interval", time.Millisecond*100, "interval between calls")
	duration = flag.Duration("duration", time.Second*5, "sustained duration")
	multiplier = flag.Float64("multiplier", 1.0, "concurrency multiplier (w.r.t max-open)")
	flag.Parse()

	log.Printf("workload= %#v; maxOpen= %#v; maxIdle= %#v; interval= %#v; duration= %#v; multiplier= %#v\n",
		*workload,
		*maxOpen,
		*maxIdle,
		*interval,
		*duration,
		*multiplier,
	)
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)

	go func() {
		sig := <-s
		panic(fmt.Sprintf("stopped; s= %#v", sig))
	}()

	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=test password=test dbname=test")
	if err != nil {
		log.Printf("cannot connect to db; err= %v\n", err)
		return
	}

	if *maxOpen != 0 || *maxIdle != 0 {
		db.SetMaxOpenConns(*maxOpen)
		db.SetMaxIdleConns(*maxIdle)
	}

	if _, err = db.Exec("TRUNCATE TABLE test"); err != nil {
		log.Printf("cannot truncate table; err= %v\n", err)
		return
	}

	go func() {
		for {
			log.Printf("[%d] db.Stats()= %#v\n", time.Now().UTC().UnixNano(), db.Stats())
			time.Sleep(time.Millisecond * 50)
		}
	}()

	wrk := worker[*workload]
	wrk(db)

	db.Close()
}

func doWork(id int, db *sql.DB) {
	log.Printf("[%d][%4d] doWork(id= %4d)\n",
		time.Now().UnixNano(),
		id,
		id,
	)
	if _, err := db.Exec(
		fmt.Sprintf("INSERT INTO test (timestamp, id) VALUES (%d, '%4d')",
			time.Now().UTC().UnixNano(),
			id,
		),
	); err != nil {
		now := time.Now().UnixNano()
		log.Printf("[%d][%4d] error inserting entry (%d, '%4d'); err= %v\n",
			now,
			id,
			now,
			id,
			err,
		)
		return
	}

	now := time.Now().UnixNano()
	log.Printf("[%d][%4d] Finished doWork()\n",
		now,
		id,
	)
}
