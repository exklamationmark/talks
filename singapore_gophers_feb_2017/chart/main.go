package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	input      *string
	output     *string
	resolution *time.Duration
	width      *int
	height     *int
)

func main() {
	input = flag.String("input", "", "input log file")
	output = flag.String("output", "", "output graph")
	resolution = flag.Duration("resolution", time.Millisecond*50, "resolution")
	width = flag.Int("width", 1000, "width")
	height = flag.Int("height", 600, "height")
	flag.Parse()

	in, err := os.Open(*input)
	if err != nil {
		log.Printf("cannot open log file; err= %v\n", err)
		return
	}
	log.Printf("reading: %s\n", *input)
	defer in.Close()

	out, err := os.Create(*output)
	if err != nil {
		log.Printf("cannot open file for writing; err= %v\n", err)
		return
	}
	log.Printf("writing: %s\n", *output)
	defer out.Close()

	data, err := extractTimeSeries(in,
		EventDoWork,
		EventPQOpen,
		EventFinished,
		EventAppendFreeConn)
	if err != nil {
		log.Printf("cannot extract time series data; err= %v\n", err)
		return
	}

	var min, max uint64
	data, min, max = normalize(data)
	// log.Printf("min= %v; max= %v\n", min, max)

	step := uint64(*resolution)
	aggrCounts, sums := aggrCount(data, min, max, step)
	log.Printf("sums= %v\n", sums)

	createChart(out, min, max, step, aggrCounts)
}
