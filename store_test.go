package httpmock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreRegister(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "store_register")
	require.Equal(t, len(mocks.mocks), 0)
	mock := New("store_register").Mock
	mocks.Register(mock)
	require.Equal(t, len(mocks.mocks), 1)
	require.Equal(t, mock.Request().Mock, mock)
	require.Equal(t, mock.Response().Mock, mock)
}

func TestStoreGetAll(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "store_get_all")
	require.Equal(t, len(mocks.mocks), 0)
	mock := New("store_get_all").Mock
	// store := mocks.GetAll()
	store := mocks
	require.Equal(t, len(mocks.mocks), 1)
	require.Equal(t, len(store.mocks), 1)
	require.Equal(t, store.mocks[0], mock)
}

func TestStoreExists(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "store_exists")
	require.Equal(t, len(mocks.mocks), 0)
	mock := New("store_exists").Mock
	require.Equal(t, len(mocks.mocks), 1)
	require.Equal(t, mocks.Exists(mock), true)
}

func TestStorePending(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "store_pending")
	New("store_pending")
	require.Equal(t, mocks.mocks, mocks.Pending())
}

func TestStoreIsPending(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "store_is_pending")
	New("store_is_pending")
	require.Equal(t, mocks.IsPending(), true)
	mocks.Flush()
	require.Equal(t, mocks.IsPending(), false)
}

func TestStoreIsDone(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "store_is_done")
	New("store_is_done")
	require.Equal(t, mocks.IsDone(), false)
	mocks.Flush()
	require.Equal(t, mocks.IsDone(), true)
}

func TestStoreRemove(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "store_remove")
	require.Equal(t, len(mocks.mocks), 0)
	mock := New("store_remove").Mock
	require.Equal(t, len(mocks.mocks), 1)
	require.Equal(t, mocks.Exists(mock), true)

	mocks.Remove(mock)
	require.Equal(t, mocks.Exists(mock), false)

	mocks.Remove(mock)
	require.Equal(t, mocks.Exists(mock), false)
}

func TestStoreFlush(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "store_flush")
	require.Equal(t, len(mocks.mocks), 0)

	mock1 := New("store_flush").Mock
	mock2 := New("store_flush").Mock
	require.Equal(t, len(mocks.mocks), 2)
	require.Equal(t, mocks.Exists(mock1), true)
	require.Equal(t, mocks.Exists(mock2), true)

	mocks.Flush()
	require.Equal(t, len(mocks.mocks), 0)
	require.Equal(t, mocks.Exists(mock1), false)
	require.Equal(t, mocks.Exists(mock2), false)
}
