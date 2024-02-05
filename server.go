package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}


// 创建一个server接口
func NewServer (ip string, port int) *Server{
	return &server {
		Ip : ip,
		Port : port
	}
}

// 启动服务器接口
func (server *Server) Start() {
	// socket listen
	net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	// accpet

	// handle

	// close listen socket
}

func main() {
	
}