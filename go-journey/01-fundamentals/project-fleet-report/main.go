package main

type Server struct {
	Name   string
	Role   string
	IP     string
	Port   int
	Online bool
	Uptime float64
}

func NewServer(name, ip, role string, port int, online bool, uptime float64) Server {
	return Server{}
}
