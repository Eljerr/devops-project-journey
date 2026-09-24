package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Server struct {
	Name   string
	Role   string
	IP     string
	Port   int
	Online bool
	Uptime float64
}

func NewServer(name, ip, role string, port int, online bool, uptime float64) Server {
	return Server{
		Name:   name,
		IP:     ip,
		Role:   role,
		Port:   port,
		Online: online,
		Uptime: uptime,
	}
}

var registry = map[string]Server{
	"web-01":      NewServer("web-01", "192.168.18.200", "web", 3000, true, 30.89),
	"monitoring":  NewServer("monitoring", "192.168.18.201", "monitoring", 30001, true, 29.32),
	"proxy":       NewServer("proxy", "192.168.18.203", "proxy", 3002, true, 49.00),
	"database-01": NewServer("database-01", "192.168.18.204", "database", 3003, true, 67.22),
	"database-02": NewServer("database-02", "192.168.18.205", "database", 3000, false, 10.89),
	"security":    NewServer("security", "192.168.18.206", "security", 3000, true, 30.89),
}

func (s Server) HealthStatus() string {
	if s.Online {
		return "UP"
	}
	return "DOWN"
}

func CheckAllServers(ctx context.Context, registry map[string]Server) []string {
	var final []string
	var wg sync.WaitGroup

	results := make(chan string, len(registry))

	for _, server := range registry {
		wg.Add(1)
		go func(n Server) {
			defer wg.Done()

			workDuration := 1 * time.Second
			if n.Name == "database-02" {
				workDuration = 3 * time.Second
			}

			select {
			case <-time.After(workDuration):
				results <- fmt.Sprintf("%s: %s", n.Name, n.HealthStatus())
			case <-ctx.Done():
				results <- fmt.Sprintf("%s: TIMEOUT (%v)", n.Name, ctx.Err())
			}
		}(server)
	}

	wg.Wait()

	for i := 0; i < len(registry); i++ {
		final = append(final, <-results)
	}

	return final
}

func main() {
	fmt.Println("=== Berikut Audit Status Server === ")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results := CheckAllServers(ctx, registry)

	for _, r := range results {
		fmt.Println(r)
	}
}
