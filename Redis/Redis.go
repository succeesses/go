package Redis

import (
	"awesomeProject/Config"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var Redis *redis.Client

func RedisInit() {

	Config.InitConfig()
	Conf := Config.ConfGlobal
	dsn1 := fmt.Sprintf("%s:%d", Conf.Redis.Ip, Conf.Redis.Port)

	rds := redis.NewClient(&redis.Options{
		Addr:     dsn1,
		Password: "",
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rds.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Redis连接失败 %v", err)
	}
	Redis = rds
}
