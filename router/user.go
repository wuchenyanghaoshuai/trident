package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/user"
)

func USER_API_ROUTER(r *gin.Engine) {
	userapi := r.Group("/api/user")
	userapi.POST("create", user.CreateUser)
	userapi.POST("delete", user.DeleteUser)
}
