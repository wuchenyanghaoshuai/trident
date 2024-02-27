package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/prometheus"
)

func Prometheus_ROUTER(r *gin.Engine) {

	cicdapi := r.Group("/prometheus")

	cicdapi.POST("getprosqlinfo", prometheus.Prometheus)

}
