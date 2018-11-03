package main

import (
	"Seckill/SecLayer/service"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var secLayerConf *service.SecLayerConf

func InitConfigT() (err error) {
	s := beego.AppConfig.String("test")
	fmt.Printf("data:%s", s)
	return
}

func InitConfig(fileType string, filePath string) (err error) {
	secLayerConf = &service.SecLayerConf{}
	conf, err := config.NewConfig(fileType, filePath)
	if err != nil {
		logs.Error("【secLayer】init config error:%v", err)
		return
	}
	//读取日志
	logPath := conf.String("logs::log_path")
	//logPath := beego.AppConfig.String("log_path")
	if len(logPath) == 0 {
		logs.Error("read log_path config error")
		err = fmt.Errorf("read log_path config error")
		return
	}
	logLevel := conf.String("logs::log_level")
	if len(logLevel) == 0 {
		logs.Error("read log_level config error")
		err = fmt.Errorf("read log_level config error")
		return
	}
	secLayerConf.LogPath = logPath
	secLayerConf.LogLevel = logLevel
	//读取redis
	redisProxy2LayerAddr := conf.String("redis::redis_proxy2layer_addr")
	if len(redisProxy2LayerAddr) == 0 {
		logs.Error("read redis_proxy2layer_addr config error")
		err = fmt.Errorf("read redis_proxy2layer_addr config error")
		return
	}
	redisProxy2LayerIdle, err := conf.Int("redis::redis_proxy2layer_idle")
	if err != nil {
		logs.Error("read redis_proxy2layer_idle config error")
		err = fmt.Errorf("read redis_proxy2layer_idle config error")
		return
	}
	redisProxy2LayerActive, err := conf.Int("redis::redis_proxy2layer_active")
	if err != nil {
		logs.Error("read redis_proxy2layer_active config error")
		err = fmt.Errorf("read redis_proxy2layer_active config error")
		return
	}
	redisProxy2LayerTimeout, err := conf.Int("redis::redis_proxy2layer_timeout")
	if err != nil {
		logs.Error("read redis_proxy2layer_idle config error")
		err = fmt.Errorf("read redis_proxy2layer_idle config error")
		return
	}
	redisProxy2LayerQueueName := conf.String("redis::redis_proxy2layer_queue_name")
	if len(redisProxy2LayerQueueName) == 0 {
		logs.Error("read redis_proxy2layer_queue_name config error")
		err = fmt.Errorf("read redis_proxy2layer_queue_name config error")
		return
	}
	secLayerConf.Proxy2LayerRedisConf.RedisAddr = redisProxy2LayerAddr
	secLayerConf.Proxy2LayerRedisConf.RedisMaxIdle = redisProxy2LayerIdle
	secLayerConf.Proxy2LayerRedisConf.RedisMaxActive = redisProxy2LayerActive
	secLayerConf.Proxy2LayerRedisConf.RedisIdleTimeout = redisProxy2LayerTimeout
	secLayerConf.Proxy2LayerRedisConf.RedisQueueName = redisProxy2LayerQueueName

	//redis:逻辑层 -> 接入层
	redisLayer2ProxyAddr := conf.String("redis::redis_layer2proxy_addr")
	if len(redisLayer2ProxyAddr) == 0 {
		logs.Error("read redis_layer2proxy_addr config error")
		err = fmt.Errorf("read redis_layer2proxy_addr config error")
		return
	}
	redisLayer2ProxyIdle, err := conf.Int("redis::redis_layer2proxy_idle")
	if err != nil {
		logs.Error("read redis_proxy2layer_idle config error")
		err = fmt.Errorf("read redis_proxy2layer_idle config error")
		return
	}
	redisLayer2ProxyActive, err := conf.Int("redis::redis_layer2proxy_active")
	if err != nil {
		logs.Error("read redis_layer2proxy_active config error")
		err = fmt.Errorf("read redis_layer2proxy_active config error")
		return
	}
	redisLayer2ProxyTimeout, err := conf.Int("redis::redis_layer2proxy_timeout")
	if err != nil {
		logs.Error("read redis_layer2proxy_timeout config error")
		err = fmt.Errorf("read redis_layer2proxy_timeout config error")
		return
	}
	redisLayer2ProxyQueueName := conf.String("redis::redis_layer2proxy_queue_name")
	if len(redisLayer2ProxyQueueName) == 0 {
		logs.Error("read redis_layer2proxy_queue_name config error")
		err = fmt.Errorf("read redis_layer2proxy_queue_name config error")
		return
	}
	secLayerConf.Layer2ProxyRedisConf.RedisAddr = redisLayer2ProxyAddr
	secLayerConf.Layer2ProxyRedisConf.RedisMaxIdle = redisLayer2ProxyIdle
	secLayerConf.Layer2ProxyRedisConf.RedisMaxActive = redisLayer2ProxyActive
	secLayerConf.Layer2ProxyRedisConf.RedisIdleTimeout = redisLayer2ProxyTimeout
	secLayerConf.Layer2ProxyRedisConf.RedisQueueName = redisLayer2ProxyQueueName

	//读取etcd
	etcdAddr := conf.String("etcd::etcd_addr")
	if len(etcdAddr) == 0 {
		logs.Error("read etcd_addr config error")
		err = fmt.Errorf("read etcd_addr config error")
		return
	}
	etcdTimeout, err := conf.Int("etcd::etcd_timeout")
	if err != nil {
		logs.Error("read etcd_timeout config error")
		err = fmt.Errorf("read etcd_timeout config error")
		return
	}
	etcdSecKillPrefix := conf.String("etcd::etcd_sec_kill_prefix")
	if len(etcdSecKillPrefix) == 0 {
		logs.Error("read etcd_sec_kill_prefix config error")
		err = fmt.Errorf("read etcd_sec_kill_prefix config error")
		return
	}
	etcdProductKey := conf.String("etcd::etcd_product_key")
	if len(etcdProductKey) == 0 {
		logs.Error("read etcd_product_key config error")
		err = fmt.Errorf("read etcd_product_key config error")
		return
	}
	secLayerConf.EtcdConfig.EtcdAddr = etcdAddr
	secLayerConf.EtcdConfig.EtcdTimeout = etcdTimeout
	secLayerConf.EtcdConfig.EtcdSecKillPrefix = etcdSecKillPrefix
	secLayerConf.EtcdConfig.EtcdProductKey = etcdProductKey

	//读取各类goroutine
	writeProxy2LayerGoroutineNum, err := conf.Int("service::write_proxy2layer_goroutine_num")
	if err != nil {
		logs.Error("read write_proxy2layer_goroutine_num config error")
		err = fmt.Errorf("read write_proxy2layer_goroutine_num config error")
		return
	}
	readLayer2ProxyGoroutineNum, err := conf.Int("service::read_layer2proxy_goroutine_num")
	if err != nil {
		logs.Error("read read_layer2proxy_goroutine_num config error")
		err = fmt.Errorf("read read_layer2proxy_goroutine_num config error")
		return
	}
	userHandleGoroutineNum, err := conf.Int("service::user_handle_goroutine_num")
	if err != nil {
		logs.Error("read user_handle_goroutine_num config error")
		err = fmt.Errorf("read user_handle_goroutine_num config error")
		return
	}
	read2handleChanSize, err := conf.Int("service::read2handle_chan_size")
	if err != nil {
		logs.Error("read read2handle_chan_size config error")
		err = fmt.Errorf("read read2handle_chan_size config error")
		return
	}
	handle2writeChanSize, err := conf.Int("service::handle2write_chan_size")
	if err != nil {
		logs.Error("read handle2write_chan_size config error")
		err = fmt.Errorf("read handle2write_chan_size config error")
		return
	}
	maxRequestTimeout, err := conf.Int("service::max_request_timeout")
	if err != nil {
		logs.Error("read max_request_timeout config error")
		err = fmt.Errorf("read max_request_timeout config error")
		return
	}
	send2writeTimeout, err := conf.Int("service::send2write_timeout")
	if err != nil {
		logs.Error("read send2write_timeout config error")
		err = fmt.Errorf("read send2write_timeout config error")
		return
	}
	send2HandleTimeout, err := conf.Int("service::send2handle_timeout")
	if err != nil {
		logs.Error("read send2handle_timeout config error")
		err = fmt.Errorf("read send2handle_timeout config error")
		return
	}
	secLayerConf.WriteGoroutineNum = writeProxy2LayerGoroutineNum
	secLayerConf.ReadGoroutineNum = readLayer2ProxyGoroutineNum
	secLayerConf.UserHandleGoroutineNum = userHandleGoroutineNum
	secLayerConf.Read2HandleChanSize = read2handleChanSize
	secLayerConf.Handle2WriteChanSize = handle2writeChanSize
	secLayerConf.MaxRequestTimeout = maxRequestTimeout
	secLayerConf.Send2HandleTimeout = send2HandleTimeout
	secLayerConf.Send2WriteTimeout = send2writeTimeout
	//加载token
	secKillTokenPasswd := conf.String("service::sec_kill_token_passwd")
	if len(secKillTokenPasswd) == 0 {
		logs.Error("read sec_kill_token_passwd config error")
		err = fmt.Errorf("read sec_kill_token_passwd config error")
		return
	}
	secLayerConf.TokenPasswd = secKillTokenPasswd
	return
}
