package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"sync"
)

type AccessLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	AccLimitLock sync.Mutex
}

func Antispam(secRequest *SecRequest) (err error) {
	_, ok := secKillConf.IdBlackMap[secRequest.UserId]
	if ok {
		logs.Error("userId:[%d] is block by black", secRequest.UserId)
		err = fmt.Errorf("invalid request")
		return
	}

	_, ok = secKillConf.IpBlackMap[secRequest.ClientAddr]
	if ok {
		logs.Error("ip:[%d] is block by black", secRequest.ClientAddr)
		err = fmt.Errorf("invalid request")
		return
	}

	limit, ok := secKillConf.AccLimitMgr.UserLimitMap[secRequest.UserId]
	if !ok {
		limit = &Limit{
			SLimit: &SecLimit{},
			MLimit: &MinLimit{},
		}
		secKillConf.AccLimitMgr.UserLimitMap[secRequest.UserId] = limit
	}

	idSecLimit := limit.SLimit.Count(secRequest.AccessTime.Unix())
	idMinLimit := limit.MLimit.Count(secRequest.AccessTime.Unix())
	if idSecLimit > secKillConf.AccLimitConf.UserSecAccessLimit {
		logs.Error("reach to idSecLimit:%d", idSecLimit)
		err = fmt.Errorf("invalid request")
		return
	}

	if idMinLimit > secKillConf.AccLimitConf.UserMinAccessLimit {
		logs.Error("reach to idMinLimit:%d", idSecLimit)
		err = fmt.Errorf("invalid request")
		return
	}

	limit, ok = secKillConf.AccLimitMgr.IpLimitMap[secRequest.ClientAddr]
	if !ok {
		limit = &Limit{
			SLimit: &SecLimit{},
			MLimit: &MinLimit{},
		}
		secKillConf.AccLimitMgr.IpLimitMap[secRequest.ClientAddr] = limit
	}

	ipSecLimit := limit.SLimit.Count(secRequest.AccessTime.Unix())
	ipMinLimit := limit.MLimit.Count(secRequest.AccessTime.Unix())
	if ipSecLimit > secKillConf.AccLimitConf.IpSecAccessLimit {
		logs.Error("reach to ipSecLimit:%d", ipSecLimit)
		err = fmt.Errorf("invalid request")
		return
	}

	if ipMinLimit > secKillConf.AccLimitConf.IpMinAccessLimit {
		logs.Error("reach to ipMinLimit:%d", ipMinLimit)
		err = fmt.Errorf("invalid request")
		return
	}
	return
}
