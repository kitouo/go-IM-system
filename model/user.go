package main

import (
	"net"
)

type User struct {
	Name string
	Address string
	Channel chan string
	Connection net.Conn

}


//创建用户
func NewUser(connection net.Conn) *User {
	userAddress := conn.RemoteAddr().String()
	user := &User {
		Name : userAddress,
		Address : userAddress,
		Channel : make(chan string),
		Connection : connection
	}

	// 启动监听当前用户channel消息的goroutine
	go user.ListenMessage()

	return user
	 

}

// 监听当前用户channel 一旦有消息立即发送给客户端
func (user *User) ListenMessage() {
	for {
		msg := <- user.Channel
		user.connection.Write([]byte(msg + "\n"))
	}
}