package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type CountryType int

const (
	CountryUS CountryType = iota
	CountrySG
)

type GDP struct {
	ID        CountryType
	GDP       sql.NullInt64
	CreatedAt time.Time
}

func main() {
	// [Slide 4: Talking to PostgreSQL, high level]
	// Goto: sql.go 570
	// It only declare the pool and start a connectionOpener
	// No real connection was created yet
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=test password=test dbname=test")
	if err != nil {
		log.Fatalf("cannot connect to db; err= %v\n", err)
	}

	// house-keeping
	// if _, err := db.Exec("TRUNCATE TABLE gdps"); err != nil {
	// 	log.Fatalf("cannot truncate data; err= %v\n", err)
	// }

	// We declare some data
	var gdps []GDP
	gdps = []GDP{
		GDP{ID: CountryUS},
		GDP{ID: CountrySG, GDP: sql.NullInt64{Valid: true, Int64: 123}},
	}

	// ---------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------
	// write
	toInsert := gdps[0]
	// Goto: sql.go 1258
	// Goto: sql.go 1241
	// Goto: sql.go 1262
	// We first get a connection (db.conn), then run a query on it
	// Back to Slide 5

	// The arguments are serialized (convert.go 44)
	// Back to Slide 6
	rows, err := db.Query(
		"INSERT INTO gdps (id, gdp) VALUES ($1, $2) RETURNING id, gdp, created_at",
		toInsert.ID,
		toInsert.GDP,
	)
	if err != nil {
		log.Fatalf("cannot insert; err= %v\n", err)
	}
	defer rows.Close()

	// ---------------------------------------------------------------------------------
	// ---------------------------------------------------------------------------------
	// read
	// Goto: sql.go 2404
	// Goto: sql.go 2111
	// what the rows actually contain is dependent on the driver
	// lib/pq conn.go 1304
	// Recall the message diagram (Slide 6)
	// Goto: sql.go 2404

	for rows.Next() {
		var inserted GDP
		if err := rows.Scan(&inserted.ID, &inserted.GDP, &inserted.CreatedAt); err != nil {
			log.Fatalf("cannot scan; err= %v\n", err)
		}
		log.Printf("%v\n", inserted)
	}
}
