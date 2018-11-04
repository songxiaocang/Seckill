package service

import "sync"

type ProductCountMgr struct {
	productCount map[int]int
	lock         sync.Mutex
}

func NewProductCountMgr() (pcountMgr *ProductCountMgr) {
	pcountMgr = &ProductCountMgr{
		productCount: make(map[int]int, 1000),
	}
	return
}

func (pcountMgr *ProductCountMgr) GetCount(productId int) (count int) {
	pcountMgr.lock.Lock()
	defer pcountMgr.lock.Unlock()
	count, _ = pcountMgr.productCount[productId]
	return
}

func (pcountMgr *ProductCountMgr) Add(productId int, count int) {
	curCount, ok := pcountMgr.productCount[productId]
	if !ok {
		curCount = count
	} else {
		curCount += count

	}
	pcountMgr.lock.Lock()
	pcountMgr.productCount[productId] = curCount
	defer pcountMgr.lock.Unlock()
}
