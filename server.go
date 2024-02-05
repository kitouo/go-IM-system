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
	return &Server {
		Ip : ip,
		Port : port,
	}
}

func (server *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	fmt.Println("连接建立成功")

}

// 启动服务器接口
func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
		
	}
	
	// close listen socket
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// handler
		go server.Handler(conn)

	}

}