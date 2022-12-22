package gock

import (
	"testing"

	"github.com/nbio/st"
)

func TestStoreRegister(t *testing.T) {
	defer after()
	mocks := register(t, "foo")
	st.Expect(t, len(mocks.mocks), 0)
	mock := New("foo").Mock
	mocks.Register(mock)
	st.Expect(t, len(mocks.mocks), 1)
	st.Expect(t, mock.Request().Mock, mock)
	st.Expect(t, mock.Response().Mock, mock)
}

func TestStoreGetAll(t *testing.T) {
	defer after()
	mocks := register(t, "foo")
	st.Expect(t, len(mocks.mocks), 0)
	mock := New("foo").Mock
	// store := mocks.GetAll()
	store := mocks
	st.Expect(t, len(mocks.mocks), 1)
	st.Expect(t, len(store.mocks), 1)
	st.Expect(t, store.mocks[0], mock)
}

func TestStoreExists(t *testing.T) {
	defer after()
	mocks := register(t, "foo")
	st.Expect(t, len(mocks.mocks), 0)
	mock := New("foo").Mock
	st.Expect(t, len(mocks.mocks), 1)
	st.Expect(t, mocks.Exists(mock), true)
}

func TestStorePending(t *testing.T) {
	defer after()
	mocks := register(t, "foo")
	New("foo")
	st.Expect(t, mocks.mocks, mocks.Pending())
}

func TestStoreIsPending(t *testing.T) {
	defer after()
	mocks := register(t, "foo")
	New("foo")
	st.Expect(t, mocks.IsPending(), true)
	mocks.Flush()
	st.Expect(t, mocks.IsPending(), false)
}

func TestStoreIsDone(t *testing.T) {
	defer after()
	mocks := register(t, "foo")
	New("foo")
	st.Expect(t, mocks.IsDone(), false)
	mocks.Flush()
	st.Expect(t, mocks.IsDone(), true)
}

func TestStoreRemove(t *testing.T) {
	defer after()
	mocks := register(t, "foo")
	st.Expect(t, len(mocks.mocks), 0)
	mock := New("foo").Mock
	st.Expect(t, len(mocks.mocks), 1)
	st.Expect(t, mocks.Exists(mock), true)

	mocks.Remove(mock)
	st.Expect(t, mocks.Exists(mock), false)

	mocks.Remove(mock)
	st.Expect(t, mocks.Exists(mock), false)
}

func TestStoreFlush(t *testing.T) {
	defer after()
	mocks := register(t, "foo")
	st.Expect(t, len(mocks.mocks), 0)

	mock1 := New("foo").Mock
	mock2 := New("foo").Mock
	st.Expect(t, len(mocks.mocks), 2)
	st.Expect(t, mocks.Exists(mock1), true)
	st.Expect(t, mocks.Exists(mock2), true)

	mocks.Flush()
	st.Expect(t, len(mocks.mocks), 0)
	st.Expect(t, mocks.Exists(mock1), false)
	st.Expect(t, mocks.Exists(mock2), false)
}
