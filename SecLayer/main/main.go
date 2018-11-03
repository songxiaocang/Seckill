package main

import (
	"Seckill/SecLayer/service"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	//初始化配置
	err := InitConfig("ini", "./conf/seclayer.conf")
	if err != nil {
		logs.Error("init config error:%v", err)
		panic(err)
		return
	}
	logs.Debug("【secLayer】init config success")
	//初始化日志
	//err = InitLogger()
	//
	//if err!=nil {
	//	logs.Error("init logger error:%v",err)
	//	panic(err)
	//	return
	//}

	logs.Debug("【secLayer】init logger success")
	//加载秒杀配置信息
	err = service.InitSecLayerConf(secLayerConf)
	if err != nil {
		logs.Error("init secLayer conf error:%v", err)
		panic(err)
		return
	}
	logs.Debug("【secLayer】init secikill config success")

	beego.Run(":8083")

}
