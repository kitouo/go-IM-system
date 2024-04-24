package main

import (
	"go-IM-system-master/domian"
)

func main() {
	server := domian.NewServer("127.0.0.1", 8080)
	server.Start()
}
