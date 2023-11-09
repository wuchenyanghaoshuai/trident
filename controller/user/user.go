package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"time"
)

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at"`
}

var db *gorm.DB

func init() {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/trident?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return
	}
	fmt.Println(db)
}

// 第一版本写的比较简单，没有把 数据库，和 用户名密码检查 和数据库加密 抽出来 下一步要抽出来让代码更简洁一些
func CreateUser(c *gin.Context) {
	var user User
	// 连接数据库

	// 绑定json数据
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{
			"message": "用户名密码输入错误",
		})
		return
	}
	// 判断用户名密码是否为空
	if user.UserName == "" || user.Password == "" {
		c.JSON(200, gin.H{
			"message": "用户名或密码不能为空",
		})
		return
	}
	db.AutoMigrate(&user)
	user.UserName = user.UserName
	user.Password = user.Password
	user.CreateAt = time.Now().Unix()
	db.WithContext(c).Table("users").Create(&user)
}

func DeleteUser(c *gin.Context) {
	var user User
	// 绑定json数据
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{
			"message": "ID输入错误",
		})
		return
	}
	//如果ins的结果为record not found 那么直接返回
	ins := db.WithContext(c).Table("users").Where("id = ?", user.Id).First(&user).Error
	if ins != nil {
		c.JSON(200, gin.H{
			"message": "没有找到匹配用户",
		})
		return
	} else {
		//db.WithContext(c).Table("users").Delete(&user)
		db.WithContext(c).Table("users").Where("id = ?", user.Id).Delete(&user)
		c.JSON(200, gin.H{
			"message": "根据id找到用户",
			"ID":      user.Id,
		})
	}
}
