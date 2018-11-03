package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func LoadProductFromEtcd(secLayerConfig *SecLayerConf) (err error) {
	etcdClient, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{secLayerConfig.EtcdConfig.EtcdAddr},
		DialTimeout: time.Duration(secLayerConfig.EtcdConfig.EtcdTimeout) * time.Second,
	})
	etcdSecKey := secLayerConfig.EtcdConfig.EtcdProductKey
	resp, err := etcdClient.Get(context.Background(), etcdSecKey)
	if err != nil {
		err = fmt.Errorf("read seckey config from etcd error:%v", err)
		return
	}

	var secProductInfoArr []SecProductInfoConf
	for _, data := range resp.Kvs {
		logs.Debug("data key:%s,data value:%s", data.Key, data.Value)

		err = json.Unmarshal(data.Value, &secProductInfoArr)
		if err != nil {
			logs.Error("unmarshal data error:%v", err)
			return
		}

		logs.Debug("unmarshal data:%v", secProductInfoArr)
	}

	updateProductInfoConf(secLayerConfig, secProductInfoArr)

	initWatchSecProductInfoConf(secLayerConfig)
	return
}

func initWatchSecProductInfoConf(secLayerConfig *SecLayerConf) {
	go watchSecProductKey(secLayerConfig)
}

func watchSecProductKey(secLayerConfig *SecLayerConf) {
	key := secLayerConfig.EtcdConfig.EtcdProductKey
	client, e := clientv3.New(clientv3.Config{
		Endpoints:   []string{secLayerConfig.EtcdConfig.EtcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if e != nil {
		logs.Error("init etcd client error:%v", e)
		return
	}

	watch := client.Watch(context.Background(), key)
	logs.Debug("start watch")
	for {
		var secProductArr []SecProductInfoConf
		var getConfSucc = true
		for wresp := range watch {
			for _, ev := range wresp.Events {
				if ev.Type.String() == mvccpb.DELETE.String() {
					logs.Debug("delete etcd config:%s", key)
					continue
				}
				if ev.Type.String() == mvccpb.PUT.String() && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &secProductArr)
					if err != nil {
						logs.Error("json unmarshal data:[%v] err :%v", ev.Kv.Value, err)
						getConfSucc = false
					}
					logs.Debug("get config from etcd succ,type:%s, key:%q,value:%q", ev.Type, ev.Kv.Key, ev.Kv.Value)
				}

				if getConfSucc {
					logs.Debug("get config from etcd succ,cofig:%v", secProductArr)
					updateProductInfoConf(secLayerConfig, secProductArr)
				}
			}
		}
	}

}

func updateProductInfoConf(secLayerConfig *SecLayerConf, secProductArr []SecProductInfoConf) {
	temp := make(map[int]*SecProductInfoConf, 1024)
	for _, data := range secProductArr {
		productInfo := data
		temp[data.ProductId] = &productInfo
	}

	secLayerContext.RwSecProductLock.Lock()
	secLayerConfig.SecProductInfoMap = temp
	secLayerContext.RwSecProductLock.Unlock()

}
