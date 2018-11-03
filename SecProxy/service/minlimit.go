package service

type TimeLimit interface {
	Count(nowTime int64) (curCount int)
	Check(nowTime int64) int
}

type MinLimit struct {
	count   int
	curTime int64
}

func (p *MinLimit) Count(nowTime int64) (count int) {
	if nowTime-p.curTime > 60 {
		p.count = 1
		p.curTime = nowTime
		count = p.count
	}
	count++
	p.count = count
	return
}

func (p *MinLimit) Check(nowTime int64) (count int) {
	if nowTime-p.curTime > 60 {
		count = 0
	}
	count = p.count
	return
}
