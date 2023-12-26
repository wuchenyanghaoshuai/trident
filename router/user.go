package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/user"
)

func USER_API_ROUTER(r *gin.Engine) {
	userapi := r.Group("/api/user")
	userapi.POST("createuser", user.CreateUser)
	userapi.POST("deleteuser", user.DeleteUser)
	userapi.POST("updateuser", user.UpdateUser)
	userapi.POST("finduser", user.FindUser)
	userapi.GET("listuser", user.ListUser)
}
