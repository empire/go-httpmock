package httpmock

import (
	"sync"
	"testing"
)

var (
	_map  = sync.Map{}
	_urls = sync.Map{}
	lock  sync.Mutex
)

func register(t *testing.T) *_mocks {
	lock.Lock()
	defer lock.Unlock()

	v := _mocks{}

	if old, ok := _map.LoadOrStore(t, &v); ok {
		return old.(*_mocks)
	}

	t.Cleanup(func() {
		_map.Delete(t)
	})

	return &v
}

func registerURL(m *_mocks, url string) {
	_urls.Store(url, m)
}

func load(url string) *_mocks {
	lock.Lock()
	defer lock.Unlock()

	m, ok := _urls.Load(url)
	if !ok {
		panic("mocks is not defined for the url")
	}
	return m.(*_mocks)
}

func IsDone(t *testing.T) bool {
	mocks, ok := _map.Load(t)
	if !ok {
		t.Errorf("TODO can't find mocks for this test")
		return false
	}
	return mocks.(*_mocks).IsDone()
}
