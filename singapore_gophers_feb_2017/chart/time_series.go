package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strconv"
)

//go:generate stringer -type=EventType
type EventType int

const (
	EventDoWork EventType = iota
	EventPQOpen
	EventFinished
	EventAppendFreeConn
	EventErrorInsert
)

var (
	regs = map[EventType]*regexp.Regexp{
		EventDoWork:         regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} \[(\d+)\]\[ *\d+\] Started doWork\(`),
		EventPQOpen:         regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} \[(\d+)\] pq\.Open\(`),
		EventFinished:       regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} \[(\d+)\]\[ *\d+\] Finished doWork\(`),
		EventAppendFreeConn: regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} \[(\d+)\] putConnDBLocked\.append\(db\.freeConn\)`),
		EventErrorInsert:    regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} \[(\d+)\]\[ *\d+\] Error inserting \(`),
	}
)

const (
	estLen = 1000
)

func extractTimeSeries(r io.Reader, events ...EventType) (map[EventType][]uint64, error) {
	res := make(map[EventType][]uint64, len(events))
	for _, e := range events {
		res[e] = make([]uint64, 0, estLen)
	}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		s := scanner.Text()
		for _, e := range events {
			if matches := regs[e].FindStringSubmatch(s); matches != nil {

				ts, err := strconv.ParseUint(matches[1], 10, 64)
				if err != nil {
					log.Printf("cannot parse timestamp (ns); str= %s; err= %v\n", matches[1], err)
					return nil, fmt.Errorf("cannot parse value s='%s' for event %v; err= %v", res[1], e, err)
				}
				res[e] = append(res[e], ts)
			}
		}
	}

	// for _, e := range events {
	// 	log.Printf("Event: %v; count= %d\n", e, len(res[e]))
	// }

	return res, nil
}

func normalize(orig map[EventType][]uint64) (_ map[EventType][]uint64, min, max uint64) {
	min = math.MaxUint64
	for _, vals := range orig {
		if len(vals) < 1 {
			continue
		}

		if vals[0] < min {
			min = vals[0]
		}
	}

	for e, vals := range orig {
		for i, ts := range vals {
			v := ts - min
			orig[e][i] = v
			if v > max {
				max = v
			}
		}
	}

	return orig, 0, max
}

func aggrCount(ts map[EventType][]uint64, min, max, step uint64) (map[EventType][]uint64, map[EventType]uint64) {
	sums := make(map[EventType]uint64, len(ts))
	tracks := make(map[EventType]uint64, len(ts))
	res := make(map[EventType][]uint64, len(ts))
	for e, _ := range ts {
		res[e] = make([]uint64, 0, (max-min)/step+1)
		tracks[e] = 0
		sums[e] = 0
	}

	var cur uint64
	for cur = step; cur <= max; cur += step {
		// move forward per event
		for e, _ := range ts {
			count := uint64(0)
			var i uint64
			for i = tracks[e]; i < uint64(len(ts[e])) && ts[e][i] < cur; i++ {
				count += 1
				sums[e] += 1
			}
			tracks[e] = i
			res[e] = append(res[e], count)
		}
	}
	for e, _ := range ts {
		count := uint64(0)
		for i := tracks[e]; i < uint64(len(ts[e])); i++ {
			count += 1
			sums[e] += 1
		}
		res[e] = append(res[e], count)
	}

	for e, _ := range ts {
		count := uint64(0)
		for i := 0; i < len(res[e]); i++ {
			count += res[e][i]
		}
	}

	return res, sums
}

// func enhance(ts map[EventType][]uint64, sums map[EventType]uint64) (map[EventType][]uint64, map[EventType]uint64) {
// 	len := len(ts[EventDoWork])
// 	ts[EventDBOpenConn] = make([]uint64, len)
//
// 	opened := 0
// 	putBack := 0
// 	for i := 0; i < len; i++ {
// 		opened += ts[EventPQOpen]
// 		putBack += ts[EventAppendFreeConn]
// 		ts[EventDBOpenConn] = ts[EventPQOpen]
// 	}
//
// 	return ts, sums
// }
