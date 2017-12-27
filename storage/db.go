package storage

import (
	"github.com/go-superman/http-watchmen/logger"
	"github.com/go-redis/redis"
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
		pong, err := client.Ping().Result()
		if err != nil {
			panic(err)
		}
		logger.Debugf("pong:%v", pong)
	}
	return client
}

