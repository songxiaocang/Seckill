package service

type Limit struct {
	SLimit TimeLimit
	MLimit TimeLimit
}

type SecLimit struct {
	count   int
	curTime int64
}

func (p *SecLimit) Count(nowTime int64) (count int) {
	if p.curTime != nowTime {
		p.count = 1
		p.curTime = nowTime
		count = p.count
	}
	count++
	p.count = count
	return
}

func (p *SecLimit) Check(nowTime int64) (count int) {
	if p.curTime != nowTime {
		count = 0
	}
	count = p.count
	return
}
