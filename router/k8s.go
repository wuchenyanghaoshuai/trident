package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/k8s/service"
	"wuchenyanghaoshuai/trident/controller/k8s/sshpod"
	"wuchenyanghaoshuai/trident/controller/k8s/storage"
	"wuchenyanghaoshuai/trident/controller/k8s/workload"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func K8S_ROUTER(r *gin.Engine) {
	k8sapi := r.Group("/api/k8s")
	r.LoadHTMLGlob("template/*")
	//ssh
	k8sapi.GET("/podssh", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ssh.html", nil)
	})

	k8sapi.GET("sshpod", sshpod.TerminalPod)
	//node
	k8sapi.GET("getnodes", workload.GetNode)
	//pods
	k8sapi.POST("getpods", workload.ListPods)
	k8sapi.POST("deletepods", workload.DeletePod)
	k8sapi.GET("logs", workload.TailLogs)
	//namespace
	k8sapi.GET("getns", workload.ListNamespace)
	k8sapi.POST("delns", workload.DeleteNamespace)

	//deployment
	k8sapi.POST("getdeploy", workload.ListDeploy)
	k8sapi.POST("restartdeploy", workload.RestartDeploy)
	k8sapi.POST("deletedeploy", workload.DeleteDeploy)

	//statefulset
	k8sapi.POST("getsts", workload.ListSts)
	k8sapi.POST("restartsts", workload.RestartSts)
	k8sapi.POST("deletests", workload.DeleteSts)

	//service
	k8sapi.POST("getsvc", service.ListSvc)
	k8sapi.POST("deletesvc", service.DelSvc)

	//ingress
	k8sapi.POST("getingress", service.ListIngress)
	k8sapi.POST("deleteingress", service.DelIngress)

	//pv pvc sc
	k8sapi.POST("getpv", storage.ListPv)
	k8sapi.POST("getpvc", storage.ListPvc)
	k8sapi.POST("getsc", storage.ListSc)
}
