package httpmock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisteredMatchers(t *testing.T) {
	t.Parallel()

	require.Equal(t, len(MatchersHeader), 7)
	require.Equal(t, len(MatchersBody), 1)
}

func TestNewMatcher(t *testing.T) {
	t.Parallel()

	matcher := NewMatcher()
	// Funcs are not comparable, checking slice length as it's better than nothing
	// See https://golang.org/pkg/reflect/#DeepEqual
	require.Equal(t, len(matcher.Matchers), len(Matchers))
	require.Equal(t, len(matcher.Get()), len(Matchers))
}

func TestNewBasicMatcher(t *testing.T) {
	t.Parallel()

	matcher := NewBasicMatcher()
	// Funcs are not comparable, checking slice length as it's better than nothing
	// See https://golang.org/pkg/reflect/#DeepEqual
	require.Equal(t, len(matcher.Matchers), len(MatchersHeader))
	require.Equal(t, len(matcher.Get()), len(MatchersHeader))
}

func TestNewEmptyMatcher(t *testing.T) {
	t.Parallel()

	matcher := NewEmptyMatcher()
	require.Equal(t, len(matcher.Matchers), 0)
	require.Equal(t, len(matcher.Get()), 0)
}

func TestMatcherAdd(t *testing.T) {
	t.Parallel()

	matcher := NewMatcher()
	require.Equal(t, len(matcher.Matchers), len(Matchers))
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	require.Equal(t, len(matcher.Get()), len(Matchers)+1)
}

func TestMatcherSet(t *testing.T) {
	t.Parallel()

	matcher := NewMatcher()
	matchers := []MatchFunc{}
	require.Equal(t, len(matcher.Matchers), len(Matchers))
	matcher.Set(matchers)
	require.Equal(t, matcher.Matchers, matchers)
	require.Equal(t, len(matcher.Get()), 0)
}

func TestMatcherGet(t *testing.T) {
	t.Parallel()

	matcher := NewMatcher()
	matchers := []MatchFunc{}
	matcher.Set(matchers)
	require.Equal(t, matcher.Get(), matchers)
}

func TestMatcherFlush(t *testing.T) {
	t.Parallel()

	matcher := NewMatcher()
	require.Equal(t, len(matcher.Matchers), len(Matchers))
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	require.Equal(t, len(matcher.Get()), len(Matchers)+1)
	matcher.Flush()
	require.Equal(t, len(matcher.Get()), 0)
}

func TestMatcherClone(t *testing.T) {
	t.Parallel()

	matcher := DefaultMatcher.Clone()
	require.Equal(t, len(matcher.Get()), len(DefaultMatcher.Get()))
}

func TestMatcher(t *testing.T) {
	t.Parallel()

	cases := []struct {
		method  string
		url     string
		matches bool
	}{
		{"GET", "http://foo.com/bar", true},
		{"GET", "http://foo.com/baz", true},
		{"GET", "http://foo.com/foo", false},
		{"POST", "http://foo.com/bar", false},
		{"POST", "http://bar.com/bar", false},
		{"GET", "http://foo.com", false},
	}

	matcher := NewMatcher()
	matcher.Flush()
	require.Equal(t, len(matcher.Matchers), 0)

	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.Method == "GET", nil
	})
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Host == "foo.com", nil
	})
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Path == "/baz" || req.URL.Path == "/bar", nil
	})

	for _, test := range cases {
		u, _ := url.Parse(test.url)
		req := &http.Request{Method: test.method, URL: u}
		matches, err := matcher.Match(req, nil)
		require.Equal(t, err, nil)
		require.Equal(t, matches, test.matches)
	}
}

//
// func TestMatchMock(t *testing.T) {
// 	cases := []struct {
// 		method  string
// 		url     string
// 		matches bool
// 	}{
// 		{"GET", "http://foo.com/bar", true},
// 		{"GET", "http://foo.com/baz", true},
// 		{"GET", "http://foo.com/foo", false},
// 		{"POST", "http://foo.com/bar", false},
// 		{"POST", "http://bar.com/bar", false},
// 		{"GET", "http://foo.com", false},
// 	}
//
// 	matcher := DefaultMatcher
// 	matcher.Flush()
// 	require.Equal(t, len(matcher.Matchers), 0)
//
// 	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
// 		return req.Method == "GET", nil
// 	})
// 	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
// 		return req.URL.Host == "foo.com", nil
// 	})
// 	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
// 		return req.URL.Path == "/baz" || req.URL.Path == "/bar", nil
// 	})
//
// 	for _, test := range cases {
// 		Flush()
// 		mock := New(test.url).method(test.method, "").Mock
//
// 		u, _ := url.Parse(test.url)
// 		req := &http.Request{Method: test.method, URL: u}
//
// 		match, err := MatchMock(req)
// 		require.Equal(t, err, nil)
// 		if test.matches {
// 			require.Equal(t, match, mock)
// 		} else {
// 			require.Equal(t, match, nil)
// 		}
// 	}
//
// 	DefaultMatcher.Matchers = Matchers
// }
