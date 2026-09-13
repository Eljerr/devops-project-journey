package main

import (
	"fmt"
	"strings"
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
		Role:   role,
		IP:     ip,
		Port:   port,
		Online: online,
		Uptime: uptime,
	}
}

func (s Server) Describe() string {
	return fmt.Sprintf("Server: %s | IP: %v:%d | Online: %t | Uptime: %.2f jam", s.Name, s.IP, s.Port, s.Online, s.Uptime)
}

var registry = map[string]Server{
	"web-01":      NewServer("web-01", "192.168.18.200", "web", 3000, true, 30.89),
	"monitoring":  NewServer("monitoring", "192.168.18.201", "monitoring", 30001, true, 29.32),
	"proxy":       NewServer("proxy", "192.168.18.203", "proxy", 3002, true, 49.00),
	"database-01": NewServer("database-01", "192.168.18.204", "database", 3003, true, 67.22),
	"database-02": NewServer("database-02", "192.168.18.205", "database", 3000, false, 10.89),
	"security":    NewServer("security", "192.168.18.206", "security", 3000, true, 30.89),
}

func GetServer(name string, registry map[string]Server) (Server, error) {
	s, exists := registry[name]
	if !exists {
		return Server{}, fmt.Errorf("server %s tidak ditemukan", name)
	}
	return s, nil
}

func (s Server) UptimeStatus() string {
	if s.Uptime >= 168 {
		return "Excellent"
	} else if s.Uptime >= 24 {
		return "Good"
	} else {
		return "Warning"
	}
}

func ServersByRole(registry map[string]Server, role string) []Server {
	var hasil []Server
	for _, server := range registry {
		if server.Role == role {
			hasil = append(hasil, server)
		}
	}
	return hasil
}

func GenerateReport(registry map[string]Server) string {
	var sb strings.Builder

	for _, server := range registry {
		fmt.Fprintf(&sb, "%s", server.Describe())
		fmt.Fprintf(&sb, "\nStatus: %s\n", server.UptimeStatus())
	}
	return sb.String()
}

func main() {
	report := GenerateReport(registry)
	fmt.Println(report)

	getrole := ServersByRole(registry, "database")
	fmt.Println(getrole)

	s, err := GetServer("web-01", registry)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println("Ditemukan: ", s.Name, s.IP)
	_, err2 := GetServer("web-99", registry)
	if err2 != nil {
		fmt.Println("Error : ", err2)
	}
}
