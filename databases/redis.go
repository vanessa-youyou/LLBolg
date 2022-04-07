package databases

import (
	"github.com/go-redis/redis"
<<<<<<< HEAD
	"os"
=======
"os"
>>>>>>> 19842c0245fc9a468ca1348d35cf3a85159a7d3a
)



var Redis *redis.Client

func InitRedis(){
<<<<<<< HEAD
	ip := os.Getenv("REDIS_IP")

	if ip == "" {
		ip = "127.0.0.1"
	}
	Redis = redis.NewClient(&redis.Options{
		Addr:     ip + "6379",
=======
 	ip := os.Getenv("REDIS_IP")

        if ip == "" {
                ip = "127.0.0.1"
        }

	Redis = redis.NewClient(&redis.Options{
		Addr:     ip + ":6379",
>>>>>>> 19842c0245fc9a468ca1348d35cf3a85159a7d3a
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var err error
	_, err = Redis.Ping().Result()
	if err != nil {
		panic(err)
	}
}
