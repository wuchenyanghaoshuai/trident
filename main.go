package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wuchenyanghaoshuai/trident/router"
)

func main() {
	r := router.InitRouter()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"code":    "200",
		})
	})
	fmt.Println("Running Server on http://127.0.0.1:8888")
	r.Run(":8888")
}
