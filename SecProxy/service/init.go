package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var secKillConf *SecKillConf

func InitService(serviceConf *SecKillConf) (err error) {
	logs.Debug("init service begin,data:%v", serviceConf)
	secKillConf = serviceConf

	LoadBlackList()

	InitProxy2LayerRedis()

	//initLayer2ProxyRedis()

	secKillConf.AccLimitMgr = &AccessLimitMgr{
		UserLimitMap: make(map[int]*Limit, 10000),
		IpLimitMap:   make(map[string]*Limit, 10000),
	}

	secKillConf.ReqChan = make(chan *SecRequest, 100)
	secKillConf.UserConnMap = make(map[string]chan *SecResult, 1000)

	InitRedisProgressFunc()

	return

}

func LoadBlackList() {
	secKillConf.IpBlackMap = make(map[string]bool, 10000)
	secKillConf.IdBlackMap = make(map[int]bool, 10000)

	InitBlackRedis()

	conn := secKillConf.RedisBlackPool.Get()
	defer conn.Close()
	return
	//获取ip黑名单
	reply, err := conn.Do("hgetall", "ipblacklist", time.Second)
	if err == redis.ErrNil {
		logs.Error("hgetall ipblacklist from redis error:%v", err)
		return
	}
	if err != nil {
		logs.Error("hgetall ipblacklist from redis error:%v", err)
		return
	}
	ipList, e := redis.Strings(reply, err)
	if e != nil {
		logs.Error("rdis parse error", e)
		return
	}

	for _, ip := range ipList {
		//i, i2 := strconv.Atoi(ip)
		secKillConf.IpBlackMap[ip] = true
	}

	//获取id黑名单
	reply, err = conn.Do("hgetall", "idblacklist", time.Second)
	if err == redis.ErrNil {
		logs.Error("hgetall idblacklist from redis error:%v", err)
		return
	}
	if err != nil {
		logs.Error("hgetall idblacklist from redis error:%v", err)
		return
	}
	idList, e := redis.Strings(reply, err)
	if e != nil {
		logs.Error("redis parse error", e)
		return
	}

	for _, id := range idList {
		v, err := strconv.Atoi(id)
		if err != nil {
			logs.Error("string convert error:%s", id)
			continue
		}
		secKillConf.IdBlackMap[v] = true
	}

	go SyncIpBlackList()

	go SyncIdBlackList()
}

func SyncIpBlackList() {
	var ipList []string
	lastTime := time.Now().Unix()
	for {
		conn := secKillConf.RedisBlackPool.Get()
		reply, err := conn.Do("BLPOP", "ipblacklist", time.Second)
		if err != nil {
			logs.Error("hgetall ipblacklist from redis error:%v", err)
			conn.Close()
			continue
		}
		ip, e := redis.String(reply, err)
		if e != nil {
			logs.Error("redis parse error:%v", e)
			conn.Close()
			continue
		}

		ipList = append(ipList, ip)

		curTime := time.Now().Unix()

		if len(ipList) > 100 || curTime-lastTime > 5 {
			for _, ip := range ipList {
				secKillConf.RwBlackLock.RLock()
				secKillConf.IpBlackMap[ip] = true
				secKillConf.RwBlackLock.RUnlock()
				lastTime = curTime
			}
		}

	}

}

func SyncIdBlackList() {
	for {
		conn := secKillConf.RedisBlackPool.Get()
		reply, err := conn.Do("BLPOP", "idblacklist", time.Second)
		if err != nil {
			logs.Error("hgetall idblacklist from redis error:%v", err)
			conn.Close()
			continue
		}
		id, e := redis.Int(reply, err)
		if e != nil {
			logs.Error("redis parse error:%v", e)
			conn.Close()
			continue
		}
		secKillConf.RwBlackLock.RLock()
		secKillConf.IdBlackMap[id] = true
		secKillConf.RwBlackLock.RUnlock()

	}

}

func InitProxy2LayerRedis() {
	secKillConf.RedisProxy2LayerPool = &redis.Pool{
		MaxIdle:     secKillConf.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisProxy2LayerConf.RedisAddr)
		},
	}

	conn := secKillConf.RedisProxy2LayerPool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		logs.Error("init proxy2layer redis error:%v", err)
		return
	}
}

func InitBlackRedis() {
	secKillConf.RedisBlackPool = &redis.Pool{
		MaxIdle:     secKillConf.RedisBlackConf.RedisMaxIdle,
		MaxActive:   secKillConf.RedisBlackConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisBlackConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisBlackConf.RedisAddr)
		},
	}

	conn := secKillConf.RedisBlackPool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		logs.Error("init black redis error:%v", err)
		return
	}
}

func InitRedisProgressFunc() {
	for i := 0; i < secKillConf.WriteProxy2LayerGoroutineNum; i++ {
		go WriteHandle()
	}
	for i := 0; i < secKillConf.ReadProxy2LayerGoroutineNum; i++ {
		go ReadHandle()
	}
	return
}
