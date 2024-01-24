package node

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wuchenyanghaoshuai/trident/controller/dao/mysql"
)

/* 根据id或者label或者hostname去查询机器信息
{
    "label":"devops"
}

{
    "id":1
}
*/
// 列出所有节点信息
func ListHosts(c *gin.Context) {
	var hosts []HostParams
	mysql.DB.WithContext(c).Table("host_params").Find(&hosts)
	c.JSON(200, gin.H{
		"message": "列出所有节点信息",
		"hosts":   hosts,
	})
}

// 根据id或者hostname或者label去查询机器信息
func GetHost(c *gin.Context) {
	var hosts []HostParams // 使用切片来存储可能的多个结果

	// 使用动态查询构建器
	query := mysql.DB.WithContext(c).Table("host_params")

	// 绑定 JSON 参数到结构体
	var params HostParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数输入错误",
		})
		return
	}

	// 构建查询条件
	if params.Id != 0 {
		query = query.Where("id = ?", params.Id)
	}
	if params.Hostname != "" {
		query = query.Where("hostname = ?", params.Hostname)
	}
	if params.Label != "" {
		query = query.Where("label = ?", params.Label)
	}

	// 执行查询
	res := query.Find(&hosts)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询过程中出现错误",
			"error":   res.Error.Error(),
		})
		return
	}

	if len(hosts) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有找到匹配的节点",
		})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"message": "成功找到匹配的节点",
		"hosts":   hosts,
	})
}

func DeleteHosts(c *gin.Context) {
	// 定义一个结构体来接收请求参数
	var params struct {
		Id       int    `json:"id"`
		Hostname string `json:"hostname"`
	}

	// 绑定 JSON 参数到结构体
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数输入错误",
		})
		return
	}

	// 使用动态查询构建器
	query := mysql.DB.WithContext(c).Table("host_params")

	// 构建查询条件
	if params.Id != 0 {
		query = query.Where("id = ?", params.Id)
	} else if params.Hostname != "" {
		query = query.Where("hostname = ?", params.Hostname)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "需要提供 id 或 hostname",
		})
		return
	}

	// 执行删除操作
	res := query.Delete(&HostParams{})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除过程中出现错误",
			"error":   res.Error.Error(),
		})
		return
	}

	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有找到匹配的节点，未执行删除操作",
		})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{
		"message": "节点删除成功",
	})
}

func UpdateHost(c *gin.Context) {
	var params HostParams

	// 绑定 JSON 参数到结构体
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数输入错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查 id 是否已提供
	if params.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "必须提供 id",
		})
		return
	}

	// 使用事务确保更新操作的原子性
	tx := mysql.DB.Begin()

	// 在更新之前先获取原始数据
	var original HostParams
	if err := tx.Where("id = ?", params.Id).First(&original).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"message": "没有找到对应的节点",
			"error":   err.Error(),
		})
		return
	}

	// 将传入的参数与原始数据合并，确保只更新提供了的字段
	if params.Hostname != "" {
		original.Hostname = params.Hostname
	}
	if params.Username != "" {
		original.Username = params.Username
	}
	if params.Password != "" {
		original.Password = params.Password
	}
	if params.Port != "" {
		original.Port = params.Port
	}
	if params.Ip != "" {
		original.Ip = params.Ip
	}
	if params.Status {
		original.Status = params.Status
	}
	if params.Osinfo != nil {
		original.Osinfo = params.Osinfo
	}
	if params.Label != "" {
		original.Label = params.Label
	}
	if params.Notes != "" {
		original.Notes = params.Notes
	}

	// 执行更新操作，忽略 id 字段
	if err := tx.Model(&original).Omit("id").Updates(original).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新过程中出现错误",
			"error":   err.Error(),
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "提交更新时出现错误",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{
		"message": "节点更新成功",
	})
}
