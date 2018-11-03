package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

func NewSecRequest() (secRequest *SecRequest) {
	secRequest = &SecRequest{
		ResultChan: make(chan *SecResult, 1),
	}
	return
}

func SecInfoList() (data []map[string]interface{}, err error, code int) {
	secKillConf.RwSecProductLock.RLock()
	defer secKillConf.RwSecProductLock.RUnlock()

	for _, v := range secKillConf.SecProductInfoMap {
		item, err, _ := SecInfoById(v.ProductId)
		if err != nil {
			logs.Error("get product info error:%v", err)
			continue
		}

		data = append(data, item)
	}
	return
}

func SecInfo(productId int) (data []map[string]interface{}, err error, code int) {
	secKillConf.RwSecProductLock.RLock()
	defer secKillConf.RwSecProductLock.RUnlock()
	item, err, code := SecInfoById(productId)
	if err != nil {
		logs.Error("get product info error:%v", err)
		return
	}

	logs.Debug("get all product info suc:%v", secKillConf.SecProductInfoMap)
	data = append(data, item)
	return
}

func SecInfoById(productId int) (data map[string]interface{}, err error, code int) {
	secKillConf.RwSecProductLock.RLock()
	defer secKillConf.RwSecProductLock.RUnlock()

	logs.Debug("secKillConf info:%v", secKillConf.SecProductInfoMap)

	v := secKillConf.SecProductInfoMap[productId]
	if v == nil {
		code = ErrorProductIdNotFound
		err = errors.New("product id not found")
		return
	}

	start := false
	end := false
	status := "success"

	curTime := time.Now().Unix()
	if curTime < v.StartTime {
		start = false
		end = false
		status = "activity not begin"
		code = ErrorActiviNotBegin
	}
	if curTime > v.StartTime {
		start = true
	}

	if curTime > v.EndTime {
		start = false
		end = true
		status = "activity has end"
		code = ErrorActiviHasEnd
	}

	if v.Status == ProductSatusSaleOut || v.Status == ProdcutStatusForceSaleOut {
		start = false
		end = true
		status = "product has sold out"
		code = ErrorActiviSaleOut
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start"] = start
	data["end"] = end
	data["status"] = status
	return
}

func SecKill(secRequest *SecRequest) (data map[string]interface{}, err error, code int) {
	secKillConf.RwSecProductLock.RLock()
	defer secKillConf.RwSecProductLock.RUnlock()

	//防刷
	err = Antispam(secRequest)
	if err != nil {
		code = ErrorUserServiceBusy
		//err = fmt.Errorf("invalid request")
		return
	}

	//用户验证
	err = UserCheck(secRequest)
	if err != nil {
		code = ErrorUserAuthFail
		return
	}

	//封装秒杀商品信息
	data, err, code = SecInfoById(secRequest.ProductId)

	//封装逻辑层响应信息
	userKey := fmt.Sprintf("%s_%s", secRequest.UserId, secRequest.ProductId)
	secKillConf.UserConnMap[userKey] = secRequest.ResultChan
	secKillConf.ReqChan <- secRequest

	ticker := time.NewTicker(time.Second * 10)
	defer func() {
		ticker.Stop()
		secKillConf.RwUserConnMapLock.Lock()
		delete(secKillConf.UserConnMap, userKey)
		secKillConf.RwUserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = ErrorProcessTimeout
		err = fmt.Errorf("process timeout")
		return
	case <-secRequest.CloseNotify:
		code = ErrorClientClosed
		err = fmt.Errorf("client closed")
		return
	case v := <-secKillConf.UserConnMap[userKey]:
		code = v.Code
		data["product_id"] = v.ProductId
		data["token"] = v.Token
		data["user_id"] = v.UserId
		return
	}

	return
}

func UserCheck(secRequest *SecRequest) (err error) {
	authData := fmt.Sprintf("%d%s", secRequest.UserId, secKillConf.CookieSecretKey)
	encrpAuthData := fmt.Sprintf("%s", md5.Sum([]byte(authData)))
	if secRequest.UserAuthSign != encrpAuthData {
		err = fmt.Errorf("user auth fail")
		return
	}
	return
}
