package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
	"os"
	"strconv"
	"time"
)

func CreateRedisInstance(opType, key string, durationInSeconds time.Duration, values ...string) (string, error) {

	cfg, err := ini.Load("controller/config/config.ini")
	if err != nil {
		os.Exit(1)
	}

	redis_host := cfg.Section("redis").Key("host").String()
	redis_port := cfg.Section("redis").Key("port").String()
	redis_username := cfg.Section("redis").Key("username").String()
	redis_password := cfg.Section("redis").Key("password").String()
	redis_db_str := cfg.Section("redis").Key("db").String()

	//获取到的redis的db的值是string类型的，需要转换为int类型
	redis_db, err := strconv.Atoi(redis_db_str)
	if err != nil {
		fmt.Println("redis.go redis_db_str转换为int失败", err)
		os.Exit(1)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host + ":" + redis_port,
		Username: redis_username,
		Password: redis_password,
		DB:       redis_db,
	})
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to ping redis:", err)
		return err.Error(), err
	}
	//fmt.Println("redis.go成功连接到redis: ", pong)
	switch opType {
	case "get":
		val, err := rdb.Get(ctx, key).Result()
		//		fmt.Println("redis.go现在输出的是val", val)
		if err != nil {
			return "", fmt.Errorf("redis.go获取值出错: %s", err)
		}
		return val, nil
	case "set":
		if len(values) == 0 {
			return "", fmt.Errorf("缺少值参数")
		}

		err := rdb.SetEX(ctx, key, values[0], time.Duration(durationInSeconds)).Err()

		if err != nil {
			return "", fmt.Errorf("设置值出错: %s", err)
		}
		return "", nil
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", opType)
	}
}

func GetRedisKey(key string) (string, error) {
	return CreateRedisInstance("get", key, 0)
}
