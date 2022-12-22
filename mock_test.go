package httpmock

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMock(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)
	require.Equal(t, mock.disabler.isDisabled(), false)
	require.Equal(t, len(mock.matcher.Get()), len(DefaultMatcher.Get()))

	require.Equal(t, mock.Request(), req)
	require.Equal(t, mock.Request().Mock, mock)
	require.Equal(t, mock.Response(), res)
	require.Equal(t, mock.Response().Mock, mock)
}

func TestMockDisable(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	require.Equal(t, mock.disabler.isDisabled(), false)
	mock.Disable()
	require.Equal(t, mock.disabler.isDisabled(), true)

	matches, err := mock.Match(&http.Request{})
	require.Equal(t, err, nil)
	require.Equal(t, matches, false)
}

func TestMockDone(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()

	mock := NewMock(req, res)
	require.Equal(t, mock.disabler.isDisabled(), false)
	require.Equal(t, mock.Done(), false)

	mock = NewMock(req, res)
	require.Equal(t, mock.disabler.isDisabled(), false)
	mock.Disable()
	require.Equal(t, mock.Done(), true)

	mock = NewMock(req, res)
	require.Equal(t, mock.disabler.isDisabled(), false)
	mock.request.Counter = 0
	require.Equal(t, mock.Done(), true)

	mock = NewMock(req, res)
	require.Equal(t, mock.disabler.isDisabled(), false)
	mock.request.Persisted = true
	require.Equal(t, mock.Done(), false)
}

func TestMockSetMatcher(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	require.Equal(t, len(mock.matcher.Get()), len(DefaultMatcher.Get()))
	matcher := NewMatcher()
	matcher.Flush()
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	mock.SetMatcher(matcher)
	require.Equal(t, len(mock.matcher.Get()), 1)
	require.Equal(t, mock.disabler.isDisabled(), false)

	matches, err := mock.Match(&http.Request{})
	require.Equal(t, err, nil)
	require.Equal(t, matches, true)
}

func TestMockAddMatcher(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	require.Equal(t, len(mock.matcher.Get()), len(DefaultMatcher.Get()))
	matcher := NewMatcher()
	matcher.Flush()
	mock.SetMatcher(matcher)
	mock.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	require.Equal(t, mock.disabler.isDisabled(), false)
	require.Equal(t, mock.matcher, matcher)

	matches, err := mock.Match(&http.Request{})
	require.Equal(t, err, nil)
	require.Equal(t, matches, true)
}

func TestMockMatch(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	matcher := NewMatcher()
	matcher.Flush()
	mock.SetMatcher(matcher)
	calls := 0
	mock.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		calls++
		return true, nil
	})
	mock.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		calls++
		return true, nil
	})
	require.Equal(t, mock.disabler.isDisabled(), false)
	require.Equal(t, mock.matcher, matcher)

	matches, err := mock.Match(&http.Request{})
	require.Equal(t, err, nil)
	require.Equal(t, calls, 2)
	require.Equal(t, matches, true)
}
