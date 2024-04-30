package main

import (
	"go-IM-system-master/entity"
)

func main() {
	server := entity.NewServer("127.0.0.1", 8080)
	server.Start()
}
