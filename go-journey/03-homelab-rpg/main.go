package main

import "time"

type Metrics struct {
	Name       string
	Up         bool
	CPUPercent float64
	MemPercent float64
	Uptime     time.Duration
}
