package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/user"
)

func USER_API_ROUTER(r *gin.Engine) {
	userapi := r.Group("/api/user")
	userapi.POST("createuser", user.CreateUser)
	userapi.DELETE("deleteuser", user.DeleteUser)
	userapi.PATCH("updateuser", user.UpdateUser)
	userapi.POST("finduser", user.FindUser)
	userapi.GET("listuser", user.ListUser)
}
