package service

type SecLimit struct {
	count   int
	curTime int64
}

func (p *SecLimit) Count(time int64) (count int) {
	if p.curTime != time {
		p.count = 1
		p.curTime = time
		count = p.count
	}
	p.count++
	p.curTime = time
	count = p.count
	return
}

func (p *SecLimit) Check(time int64) (count int) {
	if p.curTime != time {
		count = 0
	}
	count = p.count
	return
}
