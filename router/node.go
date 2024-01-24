package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/node/node"
	"wuchenyanghaoshuai/trident/controller/node/nodessh"
)

func HOST_ROUTER(r *gin.Engine) {
	r.GET("/host", node.ListHosts)
	r.GET("/ssh", nodessh.NodeSSH)
	hostapi := r.Group("/host")
	hostapi.POST("addhost", node.AddHost)
	hostapi.POST("delhost", node.DeleteHosts)
	hostapi.POST("gethost", node.GetHost)
	hostapi.POST("updatehost", node.UpdateHost)
}
