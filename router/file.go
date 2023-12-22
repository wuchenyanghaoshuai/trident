package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/file"
)

func FILE_ROUTER(r *gin.Engine) {
	fileapi := r.Group("/api/file")
	fileapi.POST("upload", file.UploadFile)
}
