package databases

import (
	"github.com/go-redis/redis"
)



var Redis *redis.Client


func InitRedis(){
	Redis = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var err error
	_, err = Redis.Ping().Result()
	if err != nil {
		panic(err)
	}
}
