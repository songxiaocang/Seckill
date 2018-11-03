package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func WriteHandle() {
	for {
		res := <-secKillConf.ReqChan
		conn := secKillConf.RedisProxy2LayerPool.Get()
		data, e := json.Marshal(&res)
		if e != nil {
			logs.Error("json marshal error:%v", e)
			conn.Close()
			continue
		}

		_, err := conn.Do("LPUSH", "sec_queue", string(data))
		if err != nil {
			logs.Error("set data to redis error:%v", err)
			conn.Close()
			continue
		}

		conn.Close()

	}
}

func ReadHandle() {
	for {
		conn := secKillConf.RedisProxy2LayerPool.Get()
		reply, err := conn.Do("LPOP", "recv_queue")
		if err != nil {
			logs.Error("no data get:%v", err)
			conn.Close()
			continue
		}
		data, err := redis.Bytes(reply, err)
		if err == redis.ErrNil {
			time.Sleep(time.Second)
			conn.Close()
			continue
		}
		if err != nil {
			logs.Error("redis parse error:%v", err)
			conn.Close()
			continue
		}
		var secResult *SecResult
		err = json.Unmarshal(data, secResult)
		if err != nil {
			logs.Error("json unmarshal error:%v", err)
			conn.Close()
			continue
		}
		userKey := fmt.Sprintf("%s_%s", secResult.UserId, secResult.ProductId)
		secKillConf.RwUserConnMapLock.Lock()
		resultChan, ok := secKillConf.UserConnMap[userKey]
		secKillConf.RwUserConnMapLock.Unlock()
		if !ok {
			conn.Close()
			logs.Error("no userid found:%v", userKey)
			return
		}
		resultChan <- secResult
		conn.Close()

	}
}
