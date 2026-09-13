package main

import "fmt"

type Server struct {
	Name   string
	IP     string
	Port   int
	Online bool
	Uptime float64
}

func NewServer(name string, ip string, port int, online bool, uptime float64) Server {
	return Server{
		Name:   name,
		IP:     ip,
		Port:   port,
		Online: online,
		Uptime: uptime,
	}
}

var registry = map[string]Server{
	"web-01":     NewServer("web-01", "192.168.18.200", 3000, true, 72.50),
	"db-01":      NewServer("db-01", "192.168.18.201", 5432, false, 10.0),
	"monitor-01": NewServer("monitor-01", "192.168.18.202", 8080, true, 30.00),
}

func GetServer(name string, registry map[string]Server) (Server, error) {
	s, exists := registry[name]
	if !exists {
		return Server{}, fmt.Errorf("server %s tidak ditemukan", name)
	}
	return s, nil
}

func ListOfflineServers(registry map[string]Server) []string {
	var offline []string
	for _, server := range registry {
		if !server.Online {
			offline = append(offline, server.Name)
		}
	}
	return offline
}

func main() {
	s, err := GetServer("web-01", registry)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println("Ditemukan: ", s.Name, s.IP)

	_, err2 := GetServer("web-99", registry)
	if err2 != nil {
		fmt.Println("Error", err2)
	}
	offline := ListOfflineServers(registry)
	fmt.Println("List server offline : ", offline)
}
