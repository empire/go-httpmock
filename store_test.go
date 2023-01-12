package httpmock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreRegister(t *testing.T) {
	t.Parallel()

	s := Server(t)
	mocks := load(s.URL)
	require.Equal(t, len(mocks.mocks), 0)
	mock := New(s.URL).Mock
	mocks.Register(mock)
	require.Equal(t, len(mocks.mocks), 1)
	require.Equal(t, mock.Request().Mock, mock)
	require.Equal(t, mock.Response().Mock, mock)
}

func TestStoreGetAll(t *testing.T) {
	t.Parallel()

	s := Server(t)
	mocks := load(s.URL)
	require.Equal(t, len(mocks.mocks), 0)
	mock := New(s.URL).Mock
	// store := mocks.GetAll()
	store := mocks
	require.Equal(t, len(mocks.mocks), 1)
	require.Equal(t, len(store.mocks), 1)
	require.Equal(t, store.mocks[0], mock)
}

func TestStoreExists(t *testing.T) {
	t.Parallel()

	s := Server(t)
	mocks := register(t)
	require.Equal(t, len(mocks.mocks), 0)
	mock := New(s.URL).Mock
	require.Equal(t, len(mocks.mocks), 1)
	require.Equal(t, mocks.Exists(mock), true)
}

func TestStorePending(t *testing.T) {
	t.Parallel()

	s := Server(t)
	mocks := load(s.URL)
	New(s.URL)
	require.Equal(t, mocks.mocks, mocks.Pending())
}

func TestStoreIsPending(t *testing.T) {
	t.Parallel()

	s := Server(t)
	mocks := load(s.URL)
	New(s.URL)
	require.Equal(t, mocks.IsPending(), true)
	mocks.Flush()
	require.Equal(t, mocks.IsPending(), false)
}

func TestStoreIsDone(t *testing.T) {
	t.Parallel()

	s := Server(t)
	mocks := load(s.URL)
	New(s.URL)
	require.Equal(t, mocks.IsDone(), false)
	mocks.Flush()
	require.Equal(t, mocks.IsDone(), true)
}

func TestStoreRemove(t *testing.T) {
	t.Parallel()

	s := Server(t)
	mocks := load(s.URL)
	require.Equal(t, len(mocks.mocks), 0)
	mock := New(s.URL).Mock
	require.Equal(t, len(mocks.mocks), 1)
	require.Equal(t, mocks.Exists(mock), true)

	mocks.Remove(mock)
	require.Equal(t, mocks.Exists(mock), false)

	mocks.Remove(mock)
	require.Equal(t, mocks.Exists(mock), false)
}

func TestStoreFlush(t *testing.T) {
	t.Parallel()

	s1 := Server(t)
	s2 := Server(t)
	mocks := load(s1.URL)
	require.Equal(t, len(mocks.mocks), 0)

	mock1 := New(s1.URL).Mock
	mock2 := New(s2.URL).Mock
	require.Equal(t, len(mocks.mocks), 2)
	require.Equal(t, mocks.Exists(mock1), true)
	require.Equal(t, mocks.Exists(mock2), true)

	mocks.Flush()
	require.Equal(t, len(mocks.mocks), 0)
	require.Equal(t, mocks.Exists(mock1), false)
	require.Equal(t, mocks.Exists(mock2), false)
}
