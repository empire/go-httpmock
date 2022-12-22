package gock

import (
	"sync"
	"testing"
)

var _map = sync.Map{}

func register(t *testing.T, url string) *_mocks {
	v := _mocks{}
	_map.Store(url, &v)
	t.Cleanup(func() {
		_map.Delete(url)
	})

	return &v
}

func load(url string) *_mocks {
	m, ok := _map.Load(url)
	if !ok {
		panic("mocks is not defined for the url")
	}
	return m.(*_mocks)
}
