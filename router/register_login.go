package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/user"
)

func RegisterAndLogin(r *gin.Engine) {
	r.POST("/login", user.Login)
	r.POST("/ldaplogin", user.LdapLogin)
	r.POST("/register", user.CreateUser)
}
