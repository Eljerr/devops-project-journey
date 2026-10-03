package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

type Metrics struct {
	Name       string        `json:"name"`
	UP         bool          `json:"up"`
	CPUPercent float64       `json:"cpupercent"`
	MemPercent float64       `json:"mempercent"`
	Uptime     time.Duration `json:"uptime"`
}

type Character struct {
	Name    string
	Level   int
	HP      int
	Stamina int
	Status  string
}

func BuildCharacter(m Metrics) Character {
	if !m.UP {
		return Character{Name: m.Name, Status: "K.O."}
	}

	uptime := int(m.Uptime / time.Minute)
	level := max(1, uptime/30)
	stamina := max(0, 100-int(math.Round(m.CPUPercent)))
	hp := max(0, 100-int(math.Round(m.MemPercent)))
	var status string

	health := min(hp, stamina)

	switch {
	case health <= 30:
		status = "Reeling"
	case health <= 50:
		status = "Weary"
	default:
		status = "Steady"
	}

	return Character{
		Name:    m.Name,
		Level:   level,
		HP:      hp,
		Stamina: stamina,
		Status:  status,
	}
}

func main() {
	servers := []Metrics{
		{Name: "web-01", UP: true, CPUPercent: 2.89, MemPercent: 3.02, Uptime: 120 * time.Minute},
		{Name: "db", UP: true, CPUPercent: 23.22, MemPercent: 21.22, Uptime: 180 * time.Minute},
		{Name: "monitoring", UP: true, CPUPercent: 15.29, MemPercent: 22.02, Uptime: 10 * time.Minute},
	}

	accounts := []Character{
		{Name: "Kiyotaka", Level: 5, HP: 80, Stamina: 15, Status: "Bleeding"},
	}

	for _, server := range servers {
		log.Printf("Nama Server : %s | UP : %t | CPU : %.2f%%", server.Name, server.UP, server.CPUPercent)
	}

	for _, account := range accounts {
		fmt.Printf("Name : %s | Level : %d", account.Name, account.Level)
	}
}
