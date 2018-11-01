package service

import (
	"errors"
	"time"
)

var secKillConf *SecKillConf
func SecInfo(productId int)(data map[string]interface{},err error,code int){
	secKillConf.RwSecProductLock.Lock()
	defer secKillConf.RwSecProductLock.Unlock()

	v,ok := secKillConf.SecProductInfoMap[productId]
	if !ok {
		code = ErrorProductIdNotFound
		err = errors.New("product id not found")
		return
	}

	start := false
	end := false
	status := "success"

	curTime := time.Now().Unix()
	if  curTime< v.StartTime{
		start = false
		end = false
		status = "activity not begin"
		code = ErrorActiviNotBegin
	}
	if curTime< v.StartTime {
		start = true
	}

	if  curTime> v.EndTime {
		start = false
		end = true
		status = "activity has end"
		code = ErrorActiviHasEnd
	}

	if v.Status == ProductSatusSaleOut || v.Status == ProdcutStatusForceSaleOut{
		start = false
		end = true
		status = "product has sold out"
		code = ErrorActiviSaleOut
	}

	data["productId"] = productId
	data["start"] = start
	data["end"] = end
	data["status"] = status
	return
}
