package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/filemgr"
)

func FILE_ROUTER(r *gin.Engine) {
	fileapi := r.Group("/api/filemgr")
	fileapi.PUT("upload", filemgr.UpsloadHandler)
}
