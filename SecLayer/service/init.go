package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func InitSecLayerConf(secLayerConf *SecLayerConf) (err error) {

	err = InitRedis(secLayerConf)
	if err != nil {
		logs.Error("init redis error:%v", err)
		return
	}

	err = InitEtcd(secLayerConf)
	if err != nil {
		logs.Error("init etcd error:%v", err)
		return
	}

	err = LoadProductFromEtcd(secLayerConf)
	if err != nil {
		logs.Error("load product from etcd error:%v", err)
		return
	}

	secLayerContext.SecLayerConfig = secLayerConf
	secLayerContext.Read2HandleChan = make(chan *SecRequest, secLayerConf.Read2HandleChanSize)
	secLayerContext.Handle2WriteChan = make(chan *SecResponse, secLayerConf.Handle2WriteChanSize)

	secLayerContext.UserHistoryMap = make(map[int]*UserHistory, 1000)

	secLayerContext.ProductcountMgr = NewProductCountMgr()
	//...
	logs.Debug("【secLayer】init all success")
	return
}

func InitEtcd(secLayerConfig *SecLayerConf) (err error) {

	conn, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{secLayerConfig.EtcdConfig.EtcdAddr},
		DialTimeout: time.Duration(secLayerConfig.EtcdConfig.EtcdTimeout) * time.Second,
	})

	if err != nil {
		err = fmt.Errorf("init etcd_client error:%v", err)
		return
	}
	secLayerContext.EtcdClient = conn
	logs.Debug("init etcd succ")
	defer conn.Close()

	return
}
