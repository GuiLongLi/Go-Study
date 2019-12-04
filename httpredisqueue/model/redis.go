package model

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)


//初始化redis
func InitRedis() (*redis.Client,error) {
	var addr = viper.GetString("common.redis.host")+":"+viper.GetString("common.redis.port")
	var auth = viper.GetString("common.redis.auth")
	var db = viper.GetInt("common.redis.db")
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       db,
	})

	pong, err := cli.Ping().Result()
	if(err != nil){
		return cli,err
	}
	log.Println("redis init ",pong)
	// Output: PONG <nil>
	return cli,nil
}
