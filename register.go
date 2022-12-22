package httpmock

import (
	"sync"
	"testing"
)

var (
	_map = sync.Map{}
	lock sync.Mutex
)

func register(t *testing.T, url string) *_mocks {
	lock.Lock()
	defer lock.Unlock()

	v := _mocks{}
	_, ok := _map.LoadOrStore(url, &v)
	if ok {
		panic("value is already exists")
	}
	t.Cleanup(func() {
		_map.Delete(url)
	})

	return &v
}

func load(url string) *_mocks {
	lock.Lock()
	defer lock.Unlock()

	m, ok := _map.Load(url)
	if !ok {
		panic("mocks is not defined for the url")
	}
	return m.(*_mocks)
}
