package workload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/public"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func TailLogs(c *gin.Context) {
	//获取url参数,如果少一个参数则返回错误
	namespace := c.Query("namespace")
	podname := c.Query("podname")
	containername := c.Query("containername")
	if namespace == "" || podname == "" || containername == "" {
		c.JSON(400, gin.H{
			"error": "params error , namepsace,podname,containername can not be empty",
		})
		return
	}

	//调用k8s api
	clientset, err := public.SetKubernetesConfig()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//判断namespace是否存在
	res := public.IsNamespaceExists(namespace)
	if !res {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "namespace not exists",
		})
		return
	}
	//
	podLogOptions := corev1.PodLogOptions{
		Container: containername,
		Follow:    true,
	}
	fmt.Println(podname, namespace, containername)
	req := clientset.CoreV1().Pods(namespace).GetLogs(podname, &podLogOptions)
	// 将 HTTP 连接升级为 WebSocket 连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer ws.Close()

	// 读取日志并将其写入 WebSocket 连接
	logs, err := req.Stream(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer logs.Close()

	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := logs.Read(buf)
			if err != nil {
				break
			}
			// 检查 WebSocket 连接是否仍然打开
			if err := ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				break
			}
		}
	}()
	// 监听客户端的关闭请求，如果连接关闭，则停止读取日志
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
	logs.Close()
}
