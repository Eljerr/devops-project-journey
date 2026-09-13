package main

import "fmt"

var slaThreshold float64 = 24.0

func main() {
	severName := "web-01"
	ip := "192.168.18.200"
	port := 8080
	status := true
	uptime := 72.50

	fmt.Printf("Server: %s | IP: %v:%d | Online: %t | Uptime: %.2f jam", severName, ip, port, status, uptime)

	fmt.Println()
	if uptime >= slaThreshold {
		fmt.Printf("SLA terpenuhi")
	} else {
		fmt.Printf("SLA belum terpenuhi")
	}
}
