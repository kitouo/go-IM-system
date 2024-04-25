package domian

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播channel
	Message chan string
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// 广播
func (server *Server) BroadCast(user *User, msg string) {

	sendMsg := "[" + user.Address + "]" + user.Name + ":" + msg

	server.Message <- sendMsg

}

// 监听Messag广播消息channel的goroutine 一旦有消息立即发送给全部在线的用户
func (server *Server) MessageLinstener() {
	for {
		msg := <-server.Message
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.Channel <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	// fmt.Println("连接建立成功")

	user := NewUser(conn, server)

	user.Online()

	// 接收客户端发送的消息
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)

			if n == 0 {
				user.OffLine()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("connetion read error: ", err)
				return
			}

			// 提取用户的消息 (去除'/n')
			msg := string(buffer[:n-1])

			// 将得到的消息进行广播
			user.doMessage(msg)

		}
	}()

	// 当前handler阻塞
	select {}

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

	// start message goroutine
	go server.MessageLinstener()

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
