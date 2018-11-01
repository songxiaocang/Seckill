package service

import (
	"github.com/astaxie/beego/logs"
)

func InitService(secKillConf SecKillConf)(err error){
	logs.Debug("init service begin,data:%v",secKillConf)
	return

}