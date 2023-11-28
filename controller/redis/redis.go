package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func CreateRedisInstance(opType, key string, durationInSeconds time.Duration, values ...string) (string, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
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
