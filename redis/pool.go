package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	Pool *redis.Pool
)

func init() {
	sRedisHost := "192.168.110.83:6379"
	nDBNum := 1 //Default is 0

	Pool = newPool(sRedisHost, nDBNum)
}

func newPool(sServer string, nDBNum int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 120 * time.Second,
		MaxActive:   30, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", sServer, redis.DialDatabase(nDBNum))
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
