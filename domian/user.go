package domian

import (
	"net"
	"strings"
)

type User struct {
	Name       string
	Address    string
	Channel    chan string
	connection net.Conn
	server     *Server
}

// 创建用户
func NewUser(connection net.Conn, server *Server) *User {
	userAddress := connection.RemoteAddr().String()
	user := &User{
		Name:       userAddress,
		Address:    userAddress,
		Channel:    make(chan string),
		connection: connection,
		server:     server,
	}

	// 启动监听当前用户channel消息的goroutine
	go user.ListenMessage()

	return user

}

// 上线
func (user *User) Online() {

	// 用户上线 将用户添加到onlineMap中
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.BroadCast(user, "已上线")

}

// 下线
func (user *User) OffLine() {
	// 用户上线 将用户从onlineMap中删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.BroadCast(user, "已下线")

}

// 给当前用户对应的客户端发送消息
func (user *User) SendMsg(msg string) {
	user.connection.Write([]byte(msg))
}

// 处理用户消息
func (currentUser *User) doMessage(msg string) {
	if msg == "online" {

		// 查询在线用户
		currentUser.server.mapLock.Lock()
		for _, user := range currentUser.server.OnlineMap {
			onlineMsg := "[" + user.Address + "]" + user.Name + ":" + "在线"
			currentUser.SendMsg(onlineMsg)
		}
		currentUser.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式 rename|张三
		newName := strings.Split(msg, "|")[1]

		// 判断用户名是否存在
		_, isExist := currentUser.server.OnlineMap[newName]
		if isExist {
			currentUser.SendMsg("当前用户名已被使用\n")
		} else {
			currentUser.server.mapLock.Lock()
			delete(currentUser.server.OnlineMap, currentUser.Name)
			currentUser.server.OnlineMap[newName] = currentUser
			currentUser.server.mapLock.Unlock()

			currentUser.Name = newName
			currentUser.SendMsg("您已更新用户名:" + currentUser.Name + "\n")

		}

	} else {
		currentUser.server.BroadCast(currentUser, msg)
	}

}

// 监听当前用户channel 一旦有消息立即发送给客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.Channel
		user.connection.Write([]byte(msg + "\n"))
	}
}
