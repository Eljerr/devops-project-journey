package main 

import (
	"context"
	"fmt"
	"sync"
	"time"
	"flag"
	"log"
)

type Server struct {
	Name string 
	Status string
	Detail string
}

type Target struct {
	Name string
	URL string
}

outpulFile := flag.string("output", "", "path file output")
flag.Parse()



