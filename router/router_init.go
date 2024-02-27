package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
加载其他路由文件中的路由
*/

// 初始化其他文件中的路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	USER_API_ROUTER(r)
	RegisterAndLogin(r)
	FILE_ROUTER(r)
	K8S_ROUTER(r)
	HOST_ROUTER(r)
	CICD_ROUTER(r)
	Prometheus_ROUTER(r)
	return r
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "X-Requested-With,authorization,dept_id,app_id,role_id,domain,tenant_id,content-type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
