package middleware

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var client *redis.Client

func NewClient() {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", viper.GetString("redis.addr"), viper.GetInt("redis.port")),
		Password:     viper.GetString("redis.pass"),
		DB:           viper.GetInt("redis.db"),
		PoolSize:     viper.GetInt("redis.size"),
		MinIdleConns: viper.GetInt("redis.conns"),
	})
	pong, err := client.Ping().Result()
	fmt.Println("init redis:", pong, err)
	//Output:PONG,<nil>
}

func GetClient() *redis.Client {
	return client
}
