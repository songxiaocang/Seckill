package main

import (
	"Seckill/SecProxy/service"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"

)



var secKillConf = &service.SecKillConf{
	SecProductInfoMap: make(map[int]*service.SecProductInfoConf,1024),
}


func initConfig() (err error){
	redisAddr := beego.AppConfig.String("redis_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")
	logs.Debug("redis_addr:%s",redisAddr)
	logs.Debug("etcd_addr:%s",etcdAddr)

	secKillConf.RedisConf.RedisAddr=redisAddr
	secKillConf.EtcdConf.EtcdAddr=etcdAddr

	if len(redisAddr) == 0 || len(etcdAddr) == 0 {
		err = fmt.Errorf("init config error, redis_addr:[%s] or etcd_addr:[%s] is not set",redisAddr,etcdAddr)
		//logs.Error("init config error, redis_addr:[%s] or etcd_addr:[%s] is not set",redis_addr,etcd_addr)
		return
	}

	redisMaxIdle, e := beego.AppConfig.Int("redis_max_idle")
	if e!=nil {
		err = fmt.Errorf("read redis_max_idle config error: %v",e)
		return
	}

	redisMaxActive, e := beego.AppConfig.Int("redis_max_active")
	if e!=nil {
		err = fmt.Errorf("read redis_max_active config error: %v",e)
		return
	}

	redisIdleTimeout, e := beego.AppConfig.Int("redis_idle_timeout")
	if e!=nil {
		err = fmt.Errorf("read redis_max_idle config error: %v",e)
		return
	}

	secKillConf.RedisConf.RedisMaxIdle=redisMaxIdle
	secKillConf.RedisConf.RedisMaxActive=redisMaxActive
	secKillConf.RedisConf.RedisIdleTimeout=redisIdleTimeout

	etcdTimeout, e := beego.AppConfig.Int("etcd_timeout")
	if e!=nil {
		err = fmt.Errorf("read etcd_timeout config error: %v",e)
		return
	}
	etcdSecKeyPrefix := beego.AppConfig.String("etcd_sec_key_prefix")
	if len(etcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("read etcd_sec_key_prefix error:%s",etcdSecKeyPrefix)
		return
	}

	secKillConf.EtcdConf.EtcdTimeout=etcdTimeout
	secKillConf.EtcdConf.EtcdSecKeyPrefix=etcdSecKeyPrefix

	etcdProductKey := beego.AppConfig.String("etcd_product_key")
	if len(etcdProductKey) == 0{
		err = fmt.Errorf("read etcd_product_key error:%s",etcdProductKey)
		return
	}

	if strings.HasSuffix(etcdSecKeyPrefix,"/") == false {
		etcdSecKeyPrefix  = etcdSecKeyPrefix+"/"
	}

	etcdProductKey = fmt.Sprintf("%s%s",etcdSecKeyPrefix,etcdProductKey)

	secKillConf.EtcdConf.EtcdProductKey=etcdProductKey

	logs.Debug("etcdTimeout:%d, etcdSecKeyPrefix:%s, etcdProductKey:%s",etcdTimeout,etcdSecKeyPrefix,etcdProductKey)

	logPath := beego.AppConfig.String("log_path")
	if len(logPath) == 0 {
		err = fmt.Errorf("read log_path config error,%s",logPath)
		return
	}
	secKillConf.LogPath = logPath
	logLevel := beego.AppConfig.String("log_level")
	if len(logLevel) == 0 {
		err = fmt.Errorf("read log_level config error,%s",logLevel)
		return
	}
	secKillConf.LogLevel = logLevel

	logs.Debug("logPath:%s",logPath)
	logs.Debug("logLevel:%s",logLevel)
	return
}