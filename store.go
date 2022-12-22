package httpmock

import (
	"sync"
)

// storeMutex is used interally for store synchronization.
var storeMutex = sync.RWMutex{}

// mocks is internally used to store registered mocks.
type _mocks struct {
	mocks []Mock
}

// Register registers a new mock in the current mocks stack.
func (mocks *_mocks) Register(mock Mock) {
	if mocks.Exists(mock) {
		return
	}

	// Make ops thread safe
	storeMutex.Lock()
	defer storeMutex.Unlock()

	// TODO move it _mocks
	// Expose mock in request/response for delegation
	mock.Request().Mock = mock
	mock.Response().Mock = mock

	// Registers the mock in the global store
	mocks.mocks = append(mocks.mocks, mock)
}

// // GetAll returns the current stack of registed mocks.
// func GetAll() []Mock {
// 	storeMutex.RLock()
// 	defer storeMutex.RUnlock()
// 	return mocks
// }

// Exists checks if the given Mock is already registered.
func (mocks *_mocks) Exists(m Mock) bool {
	storeMutex.RLock()
	defer storeMutex.RUnlock()
	for _, mock := range mocks.mocks {
		if mock == m {
			return true
		}
	}
	return false
}

// Remove removes a registered mock by reference.
func (mocks *_mocks) Remove(m Mock) {
	for i, mock := range mocks.mocks {
		if mock == m {
			storeMutex.Lock()
			mocks.mocks = append(mocks.mocks[:i], mocks.mocks[i+1:]...)
			storeMutex.Unlock()
		}
	}
}

// Flush flushes the current stack of registered mocks.
func (mocks *_mocks) Flush() {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	mocks.mocks = []Mock{}
}

// Pending returns an slice of pending mocks.
func (mocks *_mocks) Pending() []Mock {
	mocks.Clean()
	storeMutex.RLock()
	defer storeMutex.RUnlock()
	return mocks.mocks
}

// IsDone returns true if all the registered mocks has been triggered successfully.
func (mocks *_mocks) IsDone() bool {
	return !mocks.IsPending()
}

// IsPending returns true if there are pending mocks.
func (mocks *_mocks) IsPending() bool {
	return len(mocks.Pending()) > 0
}

// Clean cleans the mocks store removing disabled or obsolete mocks.
func (mocks *_mocks) Clean() {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	buf := []Mock{}
	for _, mock := range mocks.mocks {
		if mock.Done() {
			continue
		}
		buf = append(buf, mock)
	}

	mocks.mocks = buf
}
