package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/k8s/workload"
)

func K8S_ROUTER(r *gin.Engine) {
	k8sapi := r.Group("/api/k8s")
	k8sapi.GET("getnodes", workload.GetNode)
	k8sapi.GET("getns", workload.ListNamespace)
	k8sapi.POST("delns", workload.DeleteNamespace)
}
