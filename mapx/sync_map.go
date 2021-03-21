package mapx

import (
	"sync"
)

var _ MapInterface = &SyncMap{}

type SyncMap struct {
	sync.Map
}

func MakeSyncMap(size int) *SyncMap {
	return &SyncMap{}
}

func (m *SyncMap) Len() (length int) {
	length = 0
	m.Range(func(key, value interface{}) (shouldContinue bool) {
		length++
		return true
	})
	return
}
