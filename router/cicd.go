package router

import (
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/controller/cicd/cicd"
	"wuchenyanghaoshuai/trident/controller/cicd/project"
)

func CICD_ROUTER(r *gin.Engine) {

	cicdapi := r.Group("/cicd")
	cicdapi.GET("listproject", project.ListProject)
	cicdapi.POST("addproject", project.AddProject)
	cicdapi.POST("deleteproject", project.DeleteProject)
	cicdapi.POST("buildimage", cicd.CICD)

}
