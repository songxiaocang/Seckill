package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"math/rand"
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

/**
业务逻辑操作
数据流转
*/

func RunProcess() {
	for i := 0; i < secLayerContext.SecLayerConfig.ReadGoroutineNum; i++ {
		secLayerContext.WaitedGroup.Add(1)
		go HandleRead()
	}
	for i := 0; i < secLayerContext.SecLayerConfig.WriteGoroutineNum; i++ {
		secLayerContext.WaitedGroup.Add(1)
		go HandleWrite()
	}
	for i := 0; i < secLayerContext.SecLayerConfig.UserHandleGoroutineNum; i++ {
		secLayerContext.WaitedGroup.Add(1)
		go HandleUser()
	}
	logs.Debug("start process")
	secLayerContext.WaitedGroup.Wait()
	logs.Debug("process end")
}

func HandleRead() {
	logs.Debug("【secLayer】handleRead goroutine begin")
	for {
		conn := secLayerContext.Proxy2LayerRedisPool.Get()
		for {
			reply, err := conn.Do("BLPOP", secLayerContext.SecLayerConfig.Proxy2LayerRedisConf.RedisQueueName, 0)
			if err != nil {
				logs.Error("no data get from redis,key:%s,err:%v", secLayerContext.SecLayerConfig.Proxy2LayerRedisConf.RedisQueueName, err)
				continue
			}
			tmp, ok := reply.([]interface{})
			if !ok || len(tmp) != 2 {
				logs.Error("pop from queue failed")
				continue
			}

			convert, ok := tmp[1].([]byte)
			if !ok {
				logs.Error("get data from tmp err:【%v】", tmp)
			}
			var secRequest *SecRequest
			err = json.Unmarshal(convert, &secRequest)
			if err != nil {
				logs.Error("json umarsal data[%v] error:%v", string(convert), err)
				continue
			}

			now := time.Now().Unix()
			if now-secRequest.AccessTime.Unix() > int64(secLayerContext.SecLayerConfig.MaxRequestTimeout) {
				logs.Error("curTime:%v is overtime")
				continue
			}
			ticker := time.NewTicker(time.Duration(secLayerContext.SecLayerConfig.Send2HandleTimeout) * time.Second)
			select {
			case secLayerContext.Read2HandleChan <- secRequest:
			case <-ticker.C:
				logs.Warn("send to handle timeout")
				break
			}

		}
		conn.Close()
	}
}

func HandleWrite() {
	logs.Debug("【secLayer】handWrite goroutine begin")
	for data := range secLayerContext.Handle2WriteChan {
		err := send2Redis(data)
		if err != nil {
			logs.Error("send2Redis error:%v", err)
			continue
		}
	}

}

func send2Redis(secResponse *SecResponse) (err error) {
	dataStr, err := json.Marshal(secResponse)
	if err != nil {
		logs.Error("json marshal data:[%v] error:%v", secResponse, err)
		return
	}
	conn := secLayerContext.Layer2ProxyRedisPool.Get()
	defer conn.Close()
	_, err = conn.Do("RPUSH", secLayerContext.SecLayerConfig.Layer2ProxyRedisConf.RedisQueueName, string(dataStr))
	if err != nil {
		logs.Error("push data[%v] to redis error:%v", dataStr, err)
		return
	}
	return
}

func HandleUser() {
	logs.Debug("【secLayer】handle user goroutine begin")
	for data := range secLayerContext.Read2HandleChan {
		secResp, err := handleSecKill(data)
		if err != nil {
			logs.Error("handle secKill error:%v", err)
			secResp.Code = ErrServiceBusy
			continue
		}

		ticker := time.NewTicker(time.Duration(secLayerContext.SecLayerConfig.Send2WriteTimeout) * time.Second)
		select {
		case secLayerContext.Handle2WriteChan <- secResp:
		case <-ticker.C:
			logs.Warn("handle2Write timeout")
			break
		}

	}
}

func handleSecKill(secReq *SecRequest) (secResp *SecResponse, err error) {
	secResp = &SecResponse{}
	secResp.UserId = secReq.UserId
	secResp.ProductId = secReq.ProductId
	secLayerContext.RwSecProductLock.RLock()
	defer secLayerContext.RwSecProductLock.RUnlock()
	secProductInfo, ok := secLayerContext.SecLayerConfig.SecProductInfoMap[secResp.ProductId]
	if !ok {
		err = fmt.Errorf("get productInfo from config error")
		secResp.Code = ErrProductNotFound
		return
	}

	if secProductInfo.Status == StatusProductSoldOut {
		secResp.Code = ErrSoldOut
		return
	}

	now := time.Now().Unix()
	soldCount := secProductInfo.Seclimit.Count(now)
	if soldCount > secProductInfo.SoldMaxLimit {
		err = fmt.Errorf("over maxlimit")
		secResp.Code = ErrRetry
		return
	}

	secLayerContext.RwUserHistoryLock.RLock()
	userHistory, ok := secLayerContext.UserHistoryMap[secResp.UserId]
	if !ok {
		userHistory = &UserHistory{
			history: make(map[int]int, 16),
		}
		secLayerContext.UserHistoryMap[secResp.UserId] = userHistory
	}
	secLayerContext.RwUserHistoryLock.RUnlock()
	if userHistory.GetProductCount(secResp.ProductId) > secProductInfo.OnePersonBuyLimit {
		err = fmt.Errorf("over user max buy limit")
		secResp.Code = ErrAlReadyBuy
		return
	}
	if secLayerContext.ProductcountMgr.GetCount(secResp.ProductId) > secProductInfo.Total {
		err = fmt.Errorf("over maxlimit")
		secResp.Code = ErrSoldOut
		secProductInfo.Status = StatusProductSoldOut
		return
	}

	rate := rand.Float64()
	logs.Debug("curRate:%v,productBuyRate:%v,count:%v,total:%v", rate, secProductInfo.BuyRate, soldCount, secProductInfo.Total)
	if rate > secProductInfo.BuyRate {
		err = fmt.Errorf("over bue rate")
		secResp.Code = ErrRetry
		return
	}

	logs.Debug("user[%s] seckill product[%s] succ", secResp.UserId, secResp.ProductId)

	secLayerContext.ProductcountMgr.Add(secResp.ProductId, 1)
	userHistory.Add(secResp.ProductId, 1)

	//组装token
	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s", secResp.UserId, secResp.ProductId, now, secLayerContext.SecLayerConfig.TokenPasswd)

	token := fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))

	secResp.Code = SecKilSuccess
	secResp.Token = token
	secResp.TokenTime = now

	return
}
