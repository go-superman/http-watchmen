package storage

import (
	"github.com/go-superman/http-watchmen/logger"
	"github.com/go-redis/redis"
	"time"
)

var client *redis.Client

func NewClient(addr string, passwd string, db int)  (*redis.Client){
	defer func() {
		if client == nil {
			panic("redis client is nil")
		}
	}()
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: passwd, // no password set
			DB:       db,  // use default DB
		})
		go func() {
			for {
				pong, err := client.Ping().Result()
				if err != nil {
					panic(err)
				}
				time.Sleep(5*time.Second)
				logger.Debugf("pong:%v", pong)
			}
		}()
	}
	return client
}


