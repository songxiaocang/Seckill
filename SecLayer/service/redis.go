package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func InitRedis(conf *SecLayerConf) (err error) {
	secLayerContext.Proxy2LayerRedisPool, err = initRedisConf(conf.Proxy2LayerRedisConf)
	if err != nil {
		logs.Error("get proxy2layer redis pool error")
		return
	}
	secLayerContext.Layer2ProxyRedisPool, err = initRedisConf(conf.Layer2ProxyRedisConf)
	if err != nil {
		logs.Error("get layer2proxy redis pool error")
		return
	}
	return
}

func initRedisConf(redisConf RedisConf) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		MaxIdle:     redisConf.RedisMaxIdle,
		MaxActive:   redisConf.RedisMaxActive,
		IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConf.RedisAddr)
		},
	}

	conn := pool.Get()

	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping  redis error: %v", err)
		return
	}

	return
}
