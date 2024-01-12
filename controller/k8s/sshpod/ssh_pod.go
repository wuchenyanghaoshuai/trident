package sshpod

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

func TerminalPod(ctx *gin.Context) {
	var r Query
	if err := ctx.ShouldBindQuery(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	// 将 HTTP 连接升级为 websocket 连接
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	// 使用 podSSH 函数处理 websocket 连接
	PodSSH(&WSClient{
		ws:     ws,
		resize: make(chan remotecommand.TerminalSize),
	}, r)
}

// websocket 升级器配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSClient 结构体，封装了 WebSocket 连接和 resize 通道，用于在 WebSocket 和 remotecommand 之间进行数据交换。
type WSClient struct {
	// WebSocket 连接对象
	ws *websocket.Conn
	// TerminalSize 类型的通道，用于传输窗口大小调整事件
	resize chan remotecommand.TerminalSize
}

// MSG 结构体，用于解析从 WebSocket 接收到的消息。
type MSG struct {
	// 消息类型字段
	MsgType string `json:"msg_type"`
	// 窗口调整消息的行数字段
	Rows uint16 `json:"rows"`
	// 窗口调整消息的列数字段
	Cols uint16 `json:"cols"`
	// 输入消息的数据字段
	Data string `json:"data"`
}

// WSClient 的 Read 方法，实现了 io.Reader 接口，从 websocket 中读取数据。
func (c *WSClient) Read(p []byte) (n int, err error) {
	// 从 WebSocket 中读取消息
	_, message, err := c.ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	var msg MSG
	if err := json.Unmarshal(message, &msg); err != nil {
		return 0, err
	}

	// 根据消息类型进行不同的处理
	switch msg.MsgType {
	// 如果是窗口调整消息
	case "resize":
		winSize := remotecommand.TerminalSize{
			// 设置宽度
			Width: msg.Cols,
			// 设置高度
			Height: msg.Rows,
		}
		// 将 TerminalSize 对象发送到 resize 通道
		c.resize <- winSize
		return 0, nil
	// 如果是输入消息
	case "input":
		copy(p, msg.Data)
		return len(msg.Data), err
	}
	return 0, nil
}

// WSClient 的 Write 方法，实现了 io.Writer 接口，将数据写入 websocket。
func (c *WSClient) Write(p []byte) (n int, err error) {
	// 将数据作为文本消息写入 WebSocket
	err = c.ws.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}

// Next WSClient 的 Next 方法，用于从 resize 通道获取下一个 TerminalSize 事件。
func (c *WSClient) Next() *remotecommand.TerminalSize {
	// 从 resize 通道读取 TerminalSize 对象
	size := <-c.resize
	return &size
}

// podSSH 函数，这是主要的 SSH 功能逻辑，使用 kubernetes client-go 的 SPDY executor 来执行远程命令。
func PodSSH(wsClient *WSClient, q Query) {

	// 使用 kubeconfig 文件初始化 kubernetes 客户端配置
	// 请注意，你应该替换 ./config 为你的 kubeconfig 文件路径

	// 根据配置创建 kubernetes 客户端
	clientSet, err := public.SetKubernetesConfig()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	restClientConfig, err := GetRestConf()
	if err != nil {
		log.Println("getKubernetesConfig:", err)
		return
	}
	// 构造一个用于执行远程命令的请求
	request := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(q.Namespace).
		Name(q.PodName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: q.ContainerName,
			Command: []string{
				q.Command,
			},
			Stdout: true,
			Stdin:  true,
			Stderr: true,
			TTY:    true,
		}, scheme.ParameterCodec)
	// 创建 SPDY executor，用于后续的 Stream 操作
	exec, err := remotecommand.NewSPDYExecutor(restClientConfig, "POST", request.URL())
	if err != nil {
		log.Fatalf("Failed to initialize executor: %v", err)
	}

	// 开始进行 Stream 操作，即通过 websocket 执行命令
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stderr:            wsClient,
		Stdout:            wsClient,
		Stdin:             wsClient,
		Tty:               true,
		TerminalSizeQueue: wsClient,
	})
	if err != nil {
		log.Fatalf("Failed to start stream: %v", err)
	}
}

// query 结构体，用于解析和验证查询参数
type Query struct {
	Namespace     string `form:"namespace" binding:"required"`
	PodName       string `form:"pod_name" binding:"required"`
	ContainerName string `form:"container_name" binding:"required"`
	Command       string `form:"command" binding:"required"`
}

// 下面这个函数主要是为了给sshpod使用的
func GetRestConf() (restConf *rest.Config, err error) {
	kubeconfigPath := "controller/config/kube_config" // kubeconfig 文件的路径
	var kubeconfig []byte

	// 从文件中读取 kubeconfig
	kubeconfig, err = ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// 生成 rest client 配置
	restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return restConf, nil
}
