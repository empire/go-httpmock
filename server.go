package httpmock

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Server(t *testing.T) *httptest.Server {
	t.Helper()

	var server *httptest.Server
	var transport *Transport
	server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, r.URL.Host, _ = strings.Cut(server.URL, "://")
		rsp, err := transport.RoundTrip(r)
		// if !assert.NoError(t, err) {
		// 	return
		// }
		if err != nil {
			_ = err
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		body, err := io.ReadAll(rsp.Body)
		defer rsp.Body.Close()
		if !assert.NoError(t, err) {
			return
		}
		h := rw.Header()
		h2 := rsp.Header
		for k, vv := range h2 {
			for _, v := range vv {
				h.Add(k, v)
			}
		}

		rw.WriteHeader(rsp.StatusCode)
		_, err = rw.Write(body)
		assert.NoError(t, err)
	}))

	mocks := register(t, server.URL)
	transport = NewTransport(mocks)

	t.Cleanup(server.Close)
	t.Cleanup(mocks.Off)

	return server
}

// Observe(DumpNoMatchersRequest)
var DumpNoMatchersRequest ObserverFunc = func(request *http.Request, mock Mock) {
	if mock != nil && mock.Response().StatusCode != http.StatusNotImplemented {
		return
	}

	bytes, _ := httputil.DumpRequestOut(request, true)
	fmt.Println(string(bytes))
	fmt.Printf("\nMatches: false\n---\n")
}
