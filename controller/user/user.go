package user

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"time"
	"wuchenyanghaoshuai/trident/controller/mysql"
)

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username" gorm:"column:username"`
	Password string `json:"password"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at"`
	UserRole string `json:"user_role"`
}

//var db *gorm.DB
//
//func init() {
//	var err error
//	dsn := "root:123456@tcp(127.0.0.1:3306)/trident?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
//	if err != nil {
//		return
//	}
//	fmt.Println(db)
//}

// 第一版本写的比较简单，没有把 数据库，和 用户名密码检查 和数据库加密 抽出来 下一步要抽出来让代码更简洁一些
func CreateUser(c *gin.Context) {
	var user User

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
	//新增密码哈希，不将明文密码存入数据库
	b, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//mysql.Init()
	mysql.DB.AutoMigrate(&user)

	user.UserName = user.UserName
	user.Password = string(b)
	user.CreateAt = time.Now().Unix()
	user.UserRole = string(1) //admin=0, user=1
	mysql.DB.WithContext(c).Table("users").Create(&user)
	c.JSON(200, gin.H{
		"message": "创建用户成功",
		"user":    user.UserName,
	})
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
	ins := mysql.DB.WithContext(c).Table("users").Where("id = ?", user.Id).First(&user).Error
	if ins != nil {
		c.JSON(200, gin.H{
			"message": "没有找到匹配用户",
		})
		return
	} else {
		//db.WithContext(c).Table("users").Delete(&user)
		mysql.DB.WithContext(c).Table("users").Where("id = ?", user.Id).Delete(&user)
		c.JSON(200, gin.H{
			"message": "找到用户并删除",
			"ID":      user.Id,
		})
	}
}

// 更新用户信息 通过id来绑定user来更新用户信息
func UpdateUser(c *gin.Context) {
	var user User
	// 绑定json数据
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{
			"message": "ID输入错误",
		})
		return
	}

	//构建一个map来存储需要更新的字段
	updateFields := make(map[string]interface{})

	// 如果用户名不为空，添加到更新字段中
	if user.UserName != "" {
		updateFields["username"] = user.UserName
	}
	// 如果用户角色不为空，添加到更新字段中
	if user.UserRole != "" {
		updateFields["user_role"] = user.UserRole
	}
	// 如果密码不为空，添加到更新字段中
	if user.Password != "" {
		// 如果更新密码，跟创建用户一样，需要加密
		b, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		updateFields["password"] = string(b)
	}

	////通过用户传入的id来查找用户，找不到的话直接返回
	ins := mysql.DB.WithContext(c).Table("users").Where("id = ?", user.Id).First(&user).Error
	if ins != nil {
		c.JSON(200, gin.H{
			"message": "没有找到匹配用户",
		})
		return
	}
	// 如果更新字段不为空，执行更新操作
	if len(updateFields) > 0 {
		updateFields["update_at"] = time.Now().Unix()
		mysql.DB.WithContext(c).Table("users").Where("id = ?", user.Id).Updates(updateFields)
		c.JSON(200, gin.H{
			"message": "找到用户并更新",
			"ID":      user.Id,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "没有提供需要更新的字段",
		})
	}
}

// 查找用户
// 如果想返回一个自定义的数据需要自己创建一个结构体这样就可以了，如果直接返回user的话，会有密码等信息暴露出来

//	type UserResponse struct {
//		ID       uint   `json:"id"`
//		Username string `json:"username"`
//	}
//
// userResponse := UserResponse{
// ID:       user.ID,
// Username: user.UserName,
// }
func FindUser(c *gin.Context) {
	//根据用户名或者id查找用户
	var user User
	// 绑定json数据
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{
			"message": "参数输入错误",
		})
		return
	}
	//判断传递过来的参数是username还是id
	if user.UserName != "" {

		ins := mysql.DB.WithContext(c).Table("users").Where("username = ?", user.UserName).First(&user)
		if ins.Error != nil {
			c.JSON(200, gin.H{
				"message": "没有找到匹配用户",
			})
		}
		c.JSON(200, gin.H{
			"message": "成功根据username找到用户",
			"data":    user,
		})

	} else if user.Id != 0 {
		ins := mysql.DB.WithContext(c).Table("users").Where("id = ?", user.Id).First(&user)
		if ins.Error != nil {
			c.JSON(200, gin.H{
				"message": "没有找到匹配用户",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "成功根据id找到用户",
			"data":    user,
		})
	}
}
