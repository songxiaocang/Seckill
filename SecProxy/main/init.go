package main

import (
	"Seckill/SecProxy/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/garyburd/redigo/redis"
	"go.etcd.io/etcd/clientv3"
	_ "golang.org/x/net/context"
	"time"
)

var(
	redisPool *redis.Pool
	etcdClient *clientv3.Client
)


func initLogger()(err error){
	config := make(map[string]interface{})
	config["filename"] = secKillConf.LogPath
	config["level"] = convertLogLevel(secKillConf.LogLevel)

	configJson, err := json.Marshal(config)
	if err!=nil {
		err = fmt.Errorf("config json marshal error: %s",err)
		return
	}

	logs.SetLogger(logs.AdapterFile,string(configJson))
	return
}

func convertLogLevel(logLevel string) int{
	switch logLevel {
	case "debug":
		return logs.LevelDebug
	case "info":
		return logs.LevelInfo
	case "warn":
		return logs.LevelWarn
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initRedis() (err error){
	redisPool = &redis.Pool{
		MaxIdle:secKillConf.RedisConf.RedisMaxIdle,
		MaxActive:secKillConf.RedisConf.RedisMaxActive,
		IdleTimeout:time.Duration(secKillConf.RedisConf.RedisIdleTimeout)*time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",secKillConf.RedisConf.RedisAddr)
		},
	 }


	conn := redisPool.Get()

	defer conn.Close()

	_, err = conn.Do("ping")
	if err !=nil {
		logs.Error("ping redis error: %v",err)
		return
	}

	return
}

func initEtcd() (err error){

	conn,err := clientv3.New(clientv3.Config{
		Endpoints:[]string{secKillConf.EtcdConf.EtcdAddr},
		DialTimeout:time.Duration(secKillConf.EtcdConf.EtcdTimeout)*time.Second,
	})

	if err!=nil {
		err=fmt.Errorf("init etcd_client error:%v",err)
		return
	}
	etcdClient = conn
	logs.Debug("init etcd succ")
	defer conn.Close()

	return
}

//初始化秒杀配置
func loadSecConf()(err error){
	//etcdSecKey := fmt.Sprintf("%s/product",secKillConf.etcdConf.etcdSecKey)
	etcdClient,_ = clientv3.New(clientv3.Config{
		Endpoints:[]string{secKillConf.EtcdConf.EtcdAddr},
		DialTimeout:time.Duration(secKillConf.EtcdConf.EtcdTimeout)*time.Second,
	})
	etcdSecKey := secKillConf.EtcdConf.EtcdProductKey
	resp, err := etcdClient.Get(context.Background(), etcdSecKey)
	if err!=nil {
		err = fmt.Errorf("read seckey config from etcd error:%v",err)
		return
	}


	var secProductInfoArr []service.SecProductInfoConf
	for _,data := range resp.Kvs{
		logs.Debug("data key:%s,data value:%s",data.Key,data.Value)

		err = json.Unmarshal(data.Value, &secProductInfoArr)
		if err!=nil {
			logs.Error("unmarshal data error:%v",err)
			return
		}

		logs.Debug("unmarshal data:%v",secProductInfoArr)
	}

	updateProductInfoConf(secProductInfoArr)

	return
}

func initWatchSecProductInfoConf(){
	go watchSecProductKey(secKillConf.EtcdConf.EtcdProductKey)
}

func watchSecProductKey(key string){
	client, e := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if e!=nil {
		logs.Error("init etcd client error:%v",e)
		return
	}

	watch := client.Watch(context.Background(), key)
	logs.Debug("start watch")
	for  {
		var secProductArr []service.SecProductInfoConf
		var getConfSucc = true
		for wresp:= range watch{
			for _,ev:= range wresp.Events{
				if ev.Type.String() == mvccpb.DELETE.String() {
					logs.Debug("delete etcd config:%s",key)
					continue
				}
				if ev.Type.String() == mvccpb.PUT.String() && string(ev.Kv.Key)==key  {
					err := json.Unmarshal(ev.Kv.Value, &secProductArr)
					if err!=nil {
						logs.Error("json unmarshal data:[%v] err :%v",ev.Kv.Value,err)
						getConfSucc = false
					}
					logs.Debug("get config from etcd succ,type:%s, key:%q,value:%q",ev.Type,ev.Kv.Key,ev.Kv.Value)
				}

				if getConfSucc {
					logs.Debug("get config from etcd succ,cofig:%v",secProductArr)
					updateProductInfoConf(secProductArr)
				}
			}
		}
	}

}

func updateProductInfoConf(secProductArr []service.SecProductInfoConf){
	 temp := make(map[int]*service.SecProductInfoConf,1024)
	for _,data := range secProductArr{
		temp[data.ProductId] = &data
	}

	secKillConf.RwSecProductLock.Lock()
	secKillConf.SecProductInfoMap = temp
	secKillConf.RwSecProductLock.Unlock()

}

func initSec()(err error){
	err = initLogger()
	if err!=nil {
		logs.Error("init logger error: %v", err)
		return
	}

	err = initRedis()
	if err!=nil {
		logs.Error("init redis error: %v", err)
		return
	}

	err = initEtcd()
	if err!=nil {
		logs.Error("init etcd error: %v", err)
		return
	}

	err = loadSecConf()
	if err!=nil {
		logs.Error("init secConf error: %v\n", err)
		return
	}

	service.InitService(*secKillConf)
	initWatchSecProductInfoConf()


	
	logs.Info("init sec success")
	return
}