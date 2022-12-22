package httpmock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransportMatch(t *testing.T) {
	t.Parallel()

	defer after()
	const uri = "http://foo.com/transport/match"
	mocks := register(t, uri)
	New(uri).Reply(204)
	u, _ := url.Parse(uri)
	req := &http.Request{URL: u}
	res, err := NewTransport(mocks).RoundTrip(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 204)
	require.Equal(t, res.Request, req)
}

func TestTransportCannotMatch(t *testing.T) {
	t.Parallel()

	defer after()
	mocks := register(t, "http://foo.com/cannot/match")
	New("http://foo.com").Reply(204)
	u, _ := url.Parse("http://127.0.0.1:1234")
	req := &http.Request{URL: u}
	_, err := NewTransport(mocks).RoundTrip(req)
	require.Equal(t, err, ErrCannotMatch)
}

//
// func TestTransportNotIntercepting(t *testing.T) {
// 	defer after()
//
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "Hello, world")
// 	}))
// 	defer ts.Close()
//
// 	New(ts.URL).Reply(200)
// 	Disable()
//
// 	u, _ := url.Parse(ts.URL)
// 	req := &http.Request{URL: u, Header: make(http.Header)}
//
// 	res, err := NewTransport().RoundTrip(req)
// 	require.Equal(t, err, nil)
// 	require.Equal(t, Intercepting(), false)
// 	require.Equal(t, res.StatusCode, 200)
// }
