package main

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"runtime"
	"time"
)

func StartStatsSender(sdc *statsd.Client) {
	go func(c *statsd.Client) {
		for {
			c.Count("num_goroutine", runtime.NumGoroutine())
			c.Gauge("num_cgo_call", runtime.NumCgoCall())

			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			c.Count("mem_alloc", m.Alloc)
			c.Count("mem_total_alloc", m.TotalAlloc)
			c.Count("mem_sys", m.Sys)
			c.Count("mem_num_gc", m.NumGC)

			time.Sleep(30 * time.Second)
		}
	}(sdc)
}