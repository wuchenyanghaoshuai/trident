package nodessh

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// 输入消息
	messageTypeInput = "input"
	// 调整窗口大小消息
	messageTypeResize = "resize"
	// 密钥认证方式
	//	authTypeKey = "key"
	// 密码认证方式
	//	authTypePwd = "pwd"
)

// websocket 连接升级
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSClient WebSocket客户端访问对象，包含WebSocket连接对象和SSH会话对象
type WSClient struct {
	// WebSocket 连接对象
	ws         *websocket.Conn
	sshSession *ssh.Session
}

// Message 用于解析从websocket接收到的json消息
type Message struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
	Text string `json:"text"`
}

// WSClient 的 Read 方法，实现了 io.Reader 接口，从 websocket 中读取数据。
func (c *WSClient) Read(p []byte) (n int, err error) {
	// 从 WebSocket 中读取消息
	_, message, err := c.ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	msg := &Message{}
	if err := json.Unmarshal(message, msg); err != nil {
		return 0, err
	}

	switch msg.Type {
	case messageTypeInput:
		// 如果是输入消息
		return copy(p, msg.Text), err
	case messageTypeResize:
		// 如果是窗口调整消息、调整窗口大小
		return 0, c.WindowChange(msg.Rows, msg.Cols)
	default:
		return 0, fmt.Errorf("invalid message type")
	}
}

// WindowChange 改变SSH Session窗口大小
func (c *WSClient) WindowChange(rows, cols int) error {
	return c.sshSession.WindowChange(rows, cols)
}

// WSClient 的 Write 方法，实现了 io.Writer 接口，将数据写入 websocket。
func (c *WSClient) Write(p []byte) (n int, err error) {
	// 将数据作为文本消息写入 WebSocket
	err = c.ws.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}

// 建立SSH Client
func sshDial(user string, ip string, port int, privatekey string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(privatekey)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}
	// 根据认证类型选择密钥或密码认证
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}
	// SSH client配置
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// 创建SSH client
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
}

// SSHHandler 处理SSH会话
func SSHHandler(wsClient *WSClient, user, ip, command string, port int) {
	// 创建SSH client

	var PrivateKey = "controller/config/id_rsa"
	sshClient, err := sshDial(user, ip, port, PrivateKey)
	if err != nil {

		fmt.Println("输入的ip没有在已经添加的主机列表中，请重新输入！")

		return
	}
	defer sshClient.Close()

	// 创建SSH session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	wsClient.sshSession = session
	// 设置终端类型及大小
	terminalModes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 24, 80, terminalModes); err != nil {
		log.Fatal(err)
	}
	// 关联对应输入、输出流
	session.Stderr = wsClient
	session.Stdout = wsClient
	session.Stdin = wsClient
	// 在远程执行命令
	if err := session.Run(command); err != nil {
		log.Fatal(err)
	}

}

// Query 查询参数
type Query struct {
	UserName string `form:"username" binding:"required"`
	IP       string `form:"ip" binding:"required"`
	Port     int    `form:"port" binding:"required"`
	Command  string `form:"command" binding:"required,oneof=sh bash"`
}

func NodeSSH(c *gin.Context) {

	var r Query
	// 绑定并校验请求参数
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		fmt.Println(err.Error())
		return
	}

	// 将 HTTP 连接升级为 websocket 连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// 开始处理 SSH 会话
	SSHHandler(&WSClient{
		ws: ws,
	}, r.UserName, r.IP, r.Command, r.Port,
	)

}
