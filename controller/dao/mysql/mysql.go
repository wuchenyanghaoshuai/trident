package mysql

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	cfg, err := ini.Load("controller/config/config.ini")
	if err != nil {
		fmt.Println("mysql配置文件读取失败", err)
		os.Exit(1)
	}

	mysql_username := cfg.Section("mysql").Key("username").String()
	mysql_password := cfg.Section("mysql").Key("password").String()
	mysql_host := cfg.Section("mysql").Key("host").String()
	mysql_port := cfg.Section("mysql").Key("port").String()
	mysql_database := cfg.Section("mysql").Key("database").String()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql_username, mysql_password, mysql_host, mysql_port, mysql_database)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return
	}
	//fmt.Println(DB)
}
