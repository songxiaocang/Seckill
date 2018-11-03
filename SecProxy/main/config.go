package main

import (
	"Seckill/SecProxy/service"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"
)

var secKillConf = &service.SecKillConf{
	SecProductInfoMap: make(map[int]*service.SecProductInfoConf, 1024),
}

func initConfig() (err error) {
	redisBlackAddr := beego.AppConfig.String("redis_black_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")
	logs.Debug("redis_addr:%s", redisBlackAddr)
	logs.Debug("etcd_addr:%s", etcdAddr)

	secKillConf.RedisBlackConf.RedisAddr = redisBlackAddr
	secKillConf.EtcdConf.EtcdAddr = etcdAddr

	if len(redisBlackAddr) == 0 || len(etcdAddr) == 0 {
		err = fmt.Errorf("init config error, redis_black_addr:[%s] or etcd_addr:[%s] is not set", redisBlackAddr, etcdAddr)
		//logs.Error("init config error, redis_addr:[%s] or etcd_addr:[%s] is not set",redis_addr,etcd_addr)
		return
	}

	redisBlackMaxIdle, e := beego.AppConfig.Int("redis_black_max_idle")
	if e != nil {
		err = fmt.Errorf("read redis_black_max_idle config error: %v", e)
		return
	}

	redisBlackMaxActive, e := beego.AppConfig.Int("redis_black_max_active")
	if e != nil {
		err = fmt.Errorf("read redis_black_max_active config error: %v", e)
		return
	}

	redisBlackIdleTimeout, e := beego.AppConfig.Int("redis_black_idle_timeout")
	if e != nil {
		err = fmt.Errorf("read redis_black_max_idle_timeout config error: %v", e)
		return
	}

	secKillConf.RedisBlackConf.RedisMaxIdle = redisBlackMaxIdle
	secKillConf.RedisBlackConf.RedisMaxActive = redisBlackMaxActive
	secKillConf.RedisBlackConf.RedisIdleTimeout = redisBlackIdleTimeout

	//读取redis proxy-layer
	redisProxy2LayerAddr := beego.AppConfig.String("redis_black_addr")
	logs.Debug("redis_proxy2layer_addr:%s", redisBlackAddr)

	secKillConf.RedisProxy2LayerConf.RedisAddr = redisProxy2LayerAddr

	if len(redisProxy2LayerAddr) == 0 {
		err = fmt.Errorf("init config error, redis_proxy2layer_addr:[%s]  is not set", redisBlackAddr)
		return
	}

	redisProxy2LayerMaxIdle, e := beego.AppConfig.Int("redis_proxy2layer_max_idle")
	if e != nil {
		err = fmt.Errorf("read redis_proxy2layer_max_idle config error: %v", e)
		return
	}

	redisProxy2LayerMaxActive, e := beego.AppConfig.Int("redis_proxy2layer_max_active")
	if e != nil {
		err = fmt.Errorf("read redis_proxy2layer_max_active config error: %v", e)
		return
	}

	redisProxy2LayerIdleTimeout, e := beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
	if e != nil {
		err = fmt.Errorf("read redis_proxy2layer_max_idle config error: %v", e)
		return
	}

	secKillConf.RedisProxy2LayerConf.RedisMaxIdle = redisProxy2LayerMaxIdle
	secKillConf.RedisProxy2LayerConf.RedisMaxActive = redisProxy2LayerMaxActive
	secKillConf.RedisProxy2LayerConf.RedisIdleTimeout = redisProxy2LayerIdleTimeout

	//读取redis layer-proxy
	redisLayer2ProxyAddr := beego.AppConfig.String("redis_layer2proxy_addr")
	logs.Debug("redis_layer2proxy_addr:%s", redisLayer2ProxyAddr)

	secKillConf.RedisLayer2ProxyConf.RedisAddr = redisLayer2ProxyAddr

	if len(redisLayer2ProxyAddr) == 0 {
		err = fmt.Errorf("init config error, redis_layer2proxy_addr:[%s] is not set", redisBlackAddr)
		//logs.Error("init config error, redis_addr:[%s] or etcd_addr:[%s] is not set",redis_addr,etcd_addr)
		return
	}

	redisLayer2ProxyMaxIdle, e := beego.AppConfig.Int("redis_layer2proxy_max_idle")
	if e != nil {
		err = fmt.Errorf("read redis_layer2proxy_max_idle config error: %v", e)
		return
	}

	redisLayer2ProxyMaxActive, e := beego.AppConfig.Int("redis_layer2proxy_max_active")
	if e != nil {
		err = fmt.Errorf("read redis_layer2proxy_max_active config error: %v", e)
		return
	}

	redisLayer2ProxyIdleTimeout, e := beego.AppConfig.Int("redis_layer2proxy_idle_timeout")
	if e != nil {
		err = fmt.Errorf("read redis_layer2proxy_idle_timeout config error: %v", e)
		return
	}

	secKillConf.RedisLayer2ProxyConf.RedisMaxIdle = redisLayer2ProxyMaxIdle
	secKillConf.RedisLayer2ProxyConf.RedisMaxActive = redisLayer2ProxyMaxActive
	secKillConf.RedisLayer2ProxyConf.RedisIdleTimeout = redisLayer2ProxyIdleTimeout

	etcdTimeout, e := beego.AppConfig.Int("etcd_timeout")
	if e != nil {
		err = fmt.Errorf("read etcd_timeout config error: %v", e)
		return
	}
	etcdSecKeyPrefix := beego.AppConfig.String("etcd_sec_key_prefix")
	if len(etcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("read etcd_sec_key_prefix error:%s", etcdSecKeyPrefix)
		return
	}

	secKillConf.EtcdConf.EtcdTimeout = etcdTimeout
	secKillConf.EtcdConf.EtcdSecKeyPrefix = etcdSecKeyPrefix

	etcdProductKey := beego.AppConfig.String("etcd_product_key")
	if len(etcdProductKey) == 0 {
		err = fmt.Errorf("read etcd_product_key error:%s", etcdProductKey)
		return
	}

	if strings.HasSuffix(etcdSecKeyPrefix, "/") == false {
		etcdSecKeyPrefix = etcdSecKeyPrefix + "/"
	}

	etcdProductKey = fmt.Sprintf("%s%s", etcdSecKeyPrefix, etcdProductKey)

	secKillConf.EtcdConf.EtcdProductKey = etcdProductKey

	logs.Debug("etcdTimeout:%d, etcdSecKeyPrefix:%s, etcdProductKey:%s", etcdTimeout, etcdSecKeyPrefix, etcdProductKey)

	logPath := beego.AppConfig.String("log_path")
	if len(logPath) == 0 {
		err = fmt.Errorf("read log_path config error,%s", logPath)
		return
	}
	secKillConf.LogPath = logPath
	logLevel := beego.AppConfig.String("log_level")
	if len(logLevel) == 0 {
		err = fmt.Errorf("read log_level config error,%s", logLevel)
		return
	}
	secKillConf.LogLevel = logLevel

	logs.Debug("logPath:%s", logPath)
	logs.Debug("logLevel:%s", logLevel)

	cookieSecretKey := beego.AppConfig.String("cookie_secretKey")
	if len(cookieSecretKey) == 0 {
		err = fmt.Errorf("read cookie_secretKey error:[%s]", cookieSecretKey)
		return
	}
	secKillConf.CookieSecretKey = cookieSecretKey

	//读取ip id限流、白名单配置
	userSecAccessLimit, e := beego.AppConfig.Int("user_sec_access_limit")
	if e != nil {
		err = fmt.Errorf("read user_sec_access_limit config error: %v", e)
		return
	}
	ipSecAccessLimit, e := beego.AppConfig.Int("ip_sec_access_limit")
	if e != nil {
		err = fmt.Errorf("read ip_sec_access_limit config error: %v", e)
		return
	}
	userMinAccessLimit, e := beego.AppConfig.Int("user_min_access_limit")
	if e != nil {
		err = fmt.Errorf("read user_min_access_limit config error: %v", e)
		return
	}
	ipMinAccessLimit, e := beego.AppConfig.Int("ip_min_access_limit")
	if e != nil {
		err = fmt.Errorf("read ip_min_access_limit config error: %v", e)
		return
	}
	referWhitelist := beego.AppConfig.String("refer_whitelist")
	if len(referWhitelist) == 0 {
		err = fmt.Errorf("read cookie_secretKey error:[%s]", cookieSecretKey)
		return
	}
	if strings.Contains(referWhitelist, ",") {
		secKillConf.ReferWhiteList = strings.Split(referWhitelist, ",")
	}
	secKillConf.AccLimitConf.UserSecAccessLimit = userSecAccessLimit
	secKillConf.AccLimitConf.IpSecAccessLimit = ipSecAccessLimit
	secKillConf.AccLimitConf.UserMinAccessLimit = userMinAccessLimit
	secKillConf.AccLimitConf.IpMinAccessLimit = ipMinAccessLimit

	//goroutine set
	writeProxy2LayerGoroutineNum, e := beego.AppConfig.Int("write_proxy2layer_goroutine_num")
	if e != nil {
		err = fmt.Errorf("read write_proxy2layer_goroutine_num config error: %v", e)
		return
	}
	readProxy2LayerGoroutineNum, e := beego.AppConfig.Int("read_proxy2layer_goroutine_num")
	if e != nil {
		err = fmt.Errorf("read read_proxy2layer_goroutine_num config error: %v", e)
		return
	}
	secKillConf.WriteProxy2LayerGoroutineNum = writeProxy2LayerGoroutineNum
	secKillConf.ReadProxy2LayerGoroutineNum = readProxy2LayerGoroutineNum

	return
}
