package project

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"wuchenyanghaoshuai/trident/controller/dao/mysql"
)

// 为什么要加projectid是因为在gitlab里获取项目到底有多少分支的时候，需要加id才能查询到
type GitlabProject struct {
	Id          int    `json:"id"`
	ProjectName string `json:"project_name"gorm:"unique"`
	ProjectID   int    `json:"project_id"gorm:"unique"`
	ProjectUrl  string `json:"project_url"`
}

func AddProject(c *gin.Context) {

	var gitlabproject GitlabProject
	err := c.ShouldBindJSON(&gitlabproject)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	mysql.DB.AutoMigrate(&gitlabproject)
	res := mysql.DB.WithContext(c).Create(&gitlabproject)
	if res.Error != nil {
		// 如果返回的错误是因为唯一性约束违反，可以返回一个特定的错误信息
		if strings.Contains(res.Error.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{"error": "项目名称重复,请尝试其他名称"})
		} else {
			// 如果是其他类型的数据库错误
			c.JSON(http.StatusInternalServerError, gin.H{"error2": res.Error.Error()})
		}
		return
	}
	c.JSON(200, gin.H{
		"message": "项目添加成功",
		"data":    gitlabproject,
	})
}

func (GitlabProject) TableName() string {
	return "gitlab_project"
}

func DeleteProject(c *gin.Context) {
	// 定义一个结构体来接收请求参数
	var params struct {
		Id          int    `json:"id"`
		ProjectName string `json:"project_name"`
	}

	// 绑定 JSON 参数到结构体
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数输入错误",
		})
		return
	}

	// 使用动态查询构建器
	query := mysql.DB.WithContext(c).Table("gitlab_project")

	// 构建查询条件
	if params.Id != 0 {
		query = query.Where("id = ?", params.Id)
	} else if params.ProjectName != "" {
		query = query.Where("project_name = ?", params.ProjectName)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "需要提供 id 或 project_name 参数",
		})
		return
	}

	// 执行删除操作
	res := query.Delete(&GitlabProject{})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除过程中出现错误",
			"error":   res.Error.Error(),
		})
		return
	}

	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有找到匹配的项目，未执行删除操作",
		})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{
		"message": "项目删除成功",
	})
}
