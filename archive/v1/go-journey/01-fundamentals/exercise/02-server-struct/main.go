package main

import "fmt"

type Server struct {
	Name   string
	IP     string
	Port   int
	Online bool
	Uptime float64
}

func NewServer(name, ip string, port int, uptime float64) Server {
	return Server{
		Name:   name,
		IP:     ip,
		Port:   port,
		Online: true,
		Uptime: uptime,
	}
}

func (s Server) IsSLACompliant(threshold float64) bool {
	return s.Uptime >= threshold
}

func (s Server) Describe() string {
	return fmt.Sprintf("Server: %s | IP: %s:%d | Online: %t | Uptime: %.2f jam", s.Name, s.IP, s.Port, s.Online, s.Uptime)
}

func main() {
	server := NewServer("web-01", "192.168.18.200", 8080, 72.50)
	fmt.Println(server.Describe())

	if server.IsSLACompliant(24.0) {
		fmt.Println("SLA terpenuhi")
	} else {
		fmt.Println("SLA belum terpenuhi")
	}
}
