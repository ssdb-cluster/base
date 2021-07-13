// Copyright 2021 The SSDB-cluster Authors
package metric

import (
	// "fmt"
	"log"
	"time"
	"testing"
)

func TestM(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m2 := NewMetric(100)

	var max int64
    for i := 0; i < 2000; i++ {
        start := time.Now()
		time.Sleep(100 * time.Microsecond)
		ts := time.Since(start)
		m2.Add(ts.Microseconds())

		if ts.Microseconds() > max {
			max = ts.Microseconds()
		}

		time.Sleep(600 * time.Microsecond)
	}

	log.Println(max, len(m2.frames), m2.CalcUntilNow())

	var count int64 = 4
	etime := time.Now().Unix() + 1
	stime := etime - count

	// m2.samples = []*Sample{
	// 	&Sample{stime: 1623818154342546, etime: 1623818155343151, count: 1085},
	// 	&Sample{stime: 1623818155343152, etime: 1623818156197455, count: 915},
	// }
	// stime = 1623818154
	// etime = 1623818157

	for i := stime; i < etime; i ++ {
		s := i * 1000 * 1000
		e := (i + 1) * 1000 * 1000 - 1
		log.Println(m2.Calc(s, e))
	}
}