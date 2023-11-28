package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"time"
	"wuchenyanghaoshuai/trident/controller/mysql"
	"wuchenyanghaoshuai/trident/controller/redis"
)

type UserLoginForm struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var userloginform UserLoginForm
	if err := c.ShouldBindJSON(&userloginform); err != nil {

		c.JSON(200, gin.H{
			"message": "错误的请求参数",
		})
		return
	}
	var user User
	fmt.Println(user.UserName, user.Password)
	ins := mysql.DB.WithContext(c).Table("users").Where("username = ?", userloginform.UserName).First(&user).Error
	if ins != nil {
		c.JSON(200, gin.H{
			"message": "用户名或密码输入错误",
		})
		return
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userloginform.Password)); err != nil {
			c.JSON(200, gin.H{
				"message": "用户名或密码输入错误",
			})
			return
		} else {
			//登陆成功以后会生成一个token，这个token会在后面的请求中用到,存储到redis里面, username:userid=ksy, value 随机生成即可,有效期两个小时，如果两个小时都在线的话，那么就会重新生成一个token替换掉原来的token，新token的有效期为7天
			userTokenKey := fmt.Sprintf("%s:%d", user.UserName, user.Id)
			userTokenValue := xid.New().String()
			durationInSeconds := time.Second * 3600
			redis.CreateRedisInstance("set", userTokenKey, durationInSeconds, userTokenValue)
			c.JSON(200, gin.H{
				"message": "登录成功",
				"token":   userTokenValue,
			})
		}
	}
}

func RefreshToken(username string, userid int) {
	userTokenKey := fmt.Sprintf("%s:%d", username, userid)
	userTokenValue := xid.New().String()
	durationInSeconds := time.Second * 604800
	redis.CreateRedisInstance("set", userTokenKey, durationInSeconds, userTokenValue)
}
