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

	k8sapi.POST("getdeploy", workload.ListDeploy)
	k8sapi.POST("restartdeploy", workload.RestartDeploy)
	k8sapi.POST("deletedeploy", workload.DeleteDeploy)

	k8sapi.POST("getsts", workload.ListSts)
	k8sapi.POST("restartsts", workload.RestartSts)
	k8sapi.POST("deletests", workload.DeleteSts)
}
