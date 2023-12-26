package filemgr

import (
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

// 上传文件到controller下的uploadfiles目录,支持同时上传n个文件
// 如果上传的是配置文件就传到config目录
// 如果上传的是普通文件，就到普通文件的目录
func UpsloadHandler(c *gin.Context) {

	form, _ := c.MultipartForm()

	for key, files := range form.File {
		var upLoadDir string
		switch key {
		case "filemgr":
			upLoadDir = "controller/uploadfiles"
		case "config":
			upLoadDir = "controller/config"
		default:
			c.JSON(400, gin.H{
				"error": "Invalid key",
			})
			return
		}
		// 创建目录，如果它不存在
		err := os.MkdirAll(upLoadDir, 0755)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		for _, file := range files {
			// 保存文件到指定的目录
			dst := filepath.Join(upLoadDir, file.Filename)
			err = c.SaveUploadedFile(file, dst)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
		}
	}

	c.JSON(200, gin.H{
		"message": "File uploaded successfully",
	})
}
