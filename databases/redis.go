package databases

import (
	"github.com/go-redis/redis"
"os"
)



var Redis *redis.Client

func InitRedis(){
 	ip := os.Getenv("REDIS_IP")

        if ip == "" {
                ip = "127.0.0.1"
        }

	Redis = redis.NewClient(&redis.Options{
		Addr:     ip + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var err error
	_, err = Redis.Ping().Result()
	if err != nil {
		panic(err)
	}
}
