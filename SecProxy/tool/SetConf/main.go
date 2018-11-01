package main

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type SecProductInfo struct {
	ProductId int
	StartTime int
	EndTime int
	Status int
	Total int
	Left int
}

const etcdKey="/zcz/secskill/product"

func main(){
	setConfToEtcd()
}

var err error

func setConfToEtcd(){
	cli,err := clientv3.New(clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout: 5*time.Second,
	})
	if err!=nil {
		logs.Error("init etcdClient error:%v",err)
		return
	}

	logs.Debug("conn succ")
	defer cli.Close()

	var secProductInfoArr []SecProductInfo

	secProductInfoArr = append(secProductInfoArr,
		SecProductInfo{
			ProductId:2,
			StartTime:1541000797,
			EndTime:1541001797,
			Status:1,
			Total:1000,
			Left:1000,
		},
	)

	secProductInfoArr = append(secProductInfoArr,
		SecProductInfo{
			ProductId:3,
			StartTime:1541000797,
			EndTime:1541001797,
			Status:1,
			Total:2000,
			Left:2000,
		},
	)

	secProductJson, err := json.Marshal(secProductInfoArr)
	logs.Debug("marshal data:%s",string(secProductJson))
	if err!=nil {
		logs.Error("secProductInfoArr json marshal err:%v",err)
		return
	}

	context1, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(context1, etcdKey, string(secProductJson))
	cancel()
	if err!=nil {
		logs.Error("put data to etcd error:%v",err)
		return
	}

	context2, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(context2, etcdKey)
	cancel()
	if err!=nil {
		logs.Error("read data from etcd error:%v",err)
		return
	}

	for _,data := range resp.Kvs{
		logs.Debug("data key:%s, data value:%s",data.Key,data.Value)
	}

}