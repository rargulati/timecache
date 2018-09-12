package timecache

import (
	"container/list"
	"sync"
	"time"
)

type TimeCache struct {
	Q     *list.List
	qLock sync.Mutex

	M     map[string]time.Time
	mLock sync.Mutex

	span time.Duration
}

func NewTimeCache(span time.Duration) *TimeCache {
	return &TimeCache{
		Q:    list.New(),
		M:    make(map[string]time.Time),
		span: span,
	}
}

func (tc *TimeCache) Add(s string) {
	tc.mLock.Lock()
	_, ok := tc.M[s]
	if ok {
		panic("putting the same entry twice not supported")
	}
	tc.mLock.Unlock()

	tc.sweep()

	tc.mLock.Lock()
	tc.M[s] = time.Now()
	tc.mLock.Unlock()

	tc.qLock.Lock()
	tc.Q.PushFront(s)
	tc.qLock.Unlock()
}

func (tc *TimeCache) sweep() {
	for {
		tc.qLock.Lock()
		back := tc.Q.Back()
		if back == nil {
			tc.qLock.Unlock()
			return
		}
		tc.qLock.Unlock()

		v := back.Value.(string)

		tc.mLock.Lock()
		t, ok := tc.M[v]
		if !ok {
			panic("inconsistent cache state")
		}
		tc.mLock.Unlock()

		if time.Since(t) > tc.span {
			tc.qLock.Lock()
			tc.Q.Remove(back)
			tc.qLock.Unlock()

			tc.mLock.Lock()
			delete(tc.M, v)
			tc.mLock.Unlock()
		}

		return
	}
}

func (tc *TimeCache) Has(s string) bool {
	tc.mLock.Lock()
	defer tc.mLock.Unlock()

	_, ok := tc.M[s]
	return ok
}
