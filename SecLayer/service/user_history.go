package service

import "sync"

type UserHistory struct {
	history map[int]int
	lock    sync.Mutex
}

func (userHistory *UserHistory) GetProductCount(productId int) (count int) {
	userHistory.lock.Lock()
	defer userHistory.lock.Unlock()
	count, _ = userHistory.history[productId]
	return
}

func (userHistory *UserHistory) Add(productId int, count int) {
	curCount, ok := userHistory.history[productId]
	if !ok {
		curCount = count
	} else {
		curCount += count

	}
	userHistory.lock.Lock()
	userHistory.history[productId] = curCount
	defer userHistory.lock.Unlock()
}
