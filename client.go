package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//	创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}
	client.conn = conn

	// 返回对象
	return client

}

// 处理服务器响应的消息
func (client *Client) DealResponse() {
	// 永久阻塞监听
	io.Copy(os.Stdout, client.conn)

}

func (client *Client) menu() bool {
	var f int

	fmt.Println("1.聊天室")
	fmt.Println("2.私聊")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&f)

	if f >= 0 && f <= 3 {
		client.flag = f
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}

}

func (client *Client) PublicChat() {
	var chatMsg string
	// 提示用户输入消息
	fmt.Println("请输入聊天内容\n输入exit退出聊天室")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {

		// 发送给服务器
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("请输入聊天内容\n输入exit退出聊天室")
		fmt.Scanln(&chatMsg)

	}

}

func (client *Client) SelectUser() {
	sendMsg := "online\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUser()
	fmt.Println("输入聊天对象的用户名\n输入exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入聊天内容\n输入exit退出聊天")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write error", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("请输入聊天内容\n输入exit退出聊天")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUser()
		fmt.Println("输入聊天对象的用户名\n输入exit退出")
		fmt.Scanln(&remoteName)

	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write error", err)
		return false
	}

	return true

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		// 根据不同模式处理不同业务
		switch client.flag {
		case 1:
			// 聊天室
			client.PublicChat()
			break
		case 2:
			// 私聊
			client.PrivateChat()
			break
		case 3:
			// 更新用户名
			client.UpdateName()
			break

		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip")
	flag.IntVar(&serverPort, "port", 8080, "server port")
}

func main() {

	// 命令行解析过程
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("服务器连接失败")
	}

	// 单独开启一个goroutine处理服务器响应消息
	go client.DealResponse()

	fmt.Println("服务器连接成功")

	// 启动客户端业务
	client.Run()

}
