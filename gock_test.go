package httpmock

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockSimple(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).Reply(201).JSON(map[string]string{"foo": "bar"})
	res, err := http.Get(s.URL)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body)[:13], `{"foo":"bar"}`)
}

//
// func TestMockOff(t *testing.T) {
//s := Server(t)
// 	New(s.URL).Reply(201).JSON(map[string]string{"foo": "bar"})
// 	mocks.Off()
// 	_, err := http.Get(s.URL)
// 	st.Reject(t, err, nil)
// }

func TestMockBodyStringResponse(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).Reply(200).BodyString("foo bar")
	res, err := http.Get(s.URL)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body), "foo bar")
}

func TestMockBodyMatch(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).BodyString("foo bar").Reply(201).BodyString("foo foo")
	res, err := http.Post(s.URL, "text/plain", bytes.NewBuffer([]byte("foo bar")))
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}

func TestMockBodyCannotMatch(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).BodyString("foo foo").Reply(201).BodyString("foo foo")
	res, err := http.Post(s.URL, "text/plain", bytes.NewBuffer([]byte("foo bar")))
	require.NoError(t, err)

	body, _ := io.ReadAll(res.Body)
	require.Equal(t, "gock: cannot match any request", string(body))
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestMockBodyMatchCompressed(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).Compression("gzip").BodyString("foo bar").Reply(201).BodyString("foo foo")

	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write([]byte("foo bar"))
	w.Close()
	req, err := http.NewRequest("POST", s.URL, &compressed)
	require.Equal(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}

// TODO rethink about this
func TestMockBodyCannotMatchCompressed(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).Compression("gzip").BodyString("foo bar").Reply(201).BodyString("foo foo")
	res, err := http.Post(s.URL, "text/plain", bytes.NewBuffer([]byte("foo bar")))
	require.NoError(t, err)
	require.Equal(t, 500, res.StatusCode)
}

func TestMockBodyMatchJSON(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).
		Post("/bar").
		JSON(map[string]string{"foo": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	res, err := http.Post(s.URL+"/bar", "application/json", bytes.NewBuffer([]byte(`{"foo":"bar"}`)))
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body)[:13], `{"bar":"foo"}`)
}

func TestMockBodyCannotMatchJSON(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).
		Post("/bar").
		JSON(map[string]string{"bar": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	res, err := http.Post(s.URL+"/bar", "application/json", bytes.NewBuffer([]byte(`{"foo":"bar"}`)))
	require.NoError(t, err)
	require.Equal(t, 500, res.StatusCode)
}

func TestMockBodyMatchCompressedJSON(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).
		Post("/bar").
		Compression("gzip").
		JSON(map[string]string{"foo": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write([]byte(`{"foo":"bar"}`))
	w.Close()
	req, err := http.NewRequest("POST", s.URL+"/bar", &compressed)
	require.Equal(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body)[:13], `{"bar":"foo"}`)
}

func TestMockBodyCannotMatchCompressedJSON(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).
		Post("/bar").
		JSON(map[string]string{"bar": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write([]byte(`{"foo":"bar"}`))
	w.Close()
	req, err := http.NewRequest("POST", s.URL+"/bar", &compressed)
	require.Equal(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, 500, res.StatusCode)
}

func TestMockMatchHeaders(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).
		MatchHeader("Content-Type", "(.*)/plain").
		Reply(200).
		BodyString("foo foo")

	res, err := http.Post(s.URL, "text/plain", bytes.NewBuffer([]byte("foo bar")))
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}

func TestMockMap(t *testing.T) {
	t.Parallel()

	defer after()

	s := Server(t)
	mock := New(s.URL)
	mock.Map(func(req *http.Request) *http.Request {
		req.URL.Host = s.URL
		return req
	})
	mock.Reply(201).JSON(map[string]string{"foo": "bar"})

	res, err := http.Get(s.URL)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"foo":"bar"}`, string(body))
}

// TODO Uncomment
// func TestMockFilter(t *testing.T) {
// 	defer after()
//
// 	s := Server(t)
// 	mock := New(s.URL)
// 	mock.Filter(func(req *http.Request) bool {
// 		return req.URL.Host == "foo.com"
// 	})
// 	mock.Reply(201).JSON(map[string]string{"foo": "bar"})
//
// 	res, err := http.Get(s.URL)
// 	require.Equal(t, err, nil)
// 	require.Equal(t, res.StatusCode, 201)
// 	body, _ := io.ReadAll(res.Body)
// 	require.Equal(t, string(body)[:13], `{"foo":"bar"}`)
// }

// TODO uncomment
// func TestMockCounterDisabled(t *testing.T) {
// 	defer after()
// 	s := Server(t)
// 	New(s.URL).Reply(204)
// 	require.Equal(t, len(*mocks), 1)
// 	res, err := http.Get(s.URL)
// 	require.Equal(t, err, nil)
// 	require.Equal(t, res.StatusCode, 204)
// 	require.Equal(t, len(*mocks), 0)
// }

//
// func TestMockEnableNetwork(t *testing.T) {
// 	defer after()
//
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "Hello, world")
// 	}))
// 	defer ts.Close()
//
// 	EnableNetworking()
// 	defer DisableNetworking()
//
//s := Server(t)
// 	New(s.URL).Reply(204)
// 	require.Equal(t, len(GetAll()), 1)
//
// 	res, err := http.Get(ts.URL)
// 	require.Equal(t, err, nil)
// 	require.Equal(t, res.StatusCode, 204)
// 	require.Equal(t, len(GetAll()), 0)
//
// 	res, err = http.Get(ts.URL)
// 	require.Equal(t, err, nil)
// 	require.Equal(t, res.StatusCode, 200)
// }
//
// func TestMockEnableNetworkFilter(t *testing.T) {
// 	defer after()
//
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "Hello, world")
// 	}))
// 	defer ts.Close()
//
// 	EnableNetworking()
// 	defer DisableNetworking()
//
// 	NetworkingFilter(func(req *http.Request) bool {
// 		return strings.Contains(req.URL.Host, "127.0.0.1")
// 	})
// 	defer DisableNetworkingFilters()
//
//s := Server(t)
// 	New(s.URL).Reply(0).SetHeader("foo", "bar")
// 	require.Equal(t, len(GetAll()), 1)
//
// 	res, err := http.Get(ts.URL)
// 	require.Equal(t, err, nil)
// 	require.Equal(t, res.StatusCode, 200)
// 	require.Equal(t, res.Header.Get("foo"), "bar")
// 	require.Equal(t, len(GetAll()), 0)
// }

func TestMockPersistent(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).
		Get("/bar").
		Persist().
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	for i := 0; i < 5; i++ {
		res, err := http.Get(s.URL + "/bar")
		require.Equal(t, err, nil)
		require.Equal(t, res.StatusCode, 200)
		body, _ := io.ReadAll(res.Body)
		require.Equal(t, string(body)[:13], `{"foo":"bar"}`)
	}
}

func TestMockPersistTimes(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).
		Get("/bar").
		Times(4).
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	for i := 0; i < 5; i++ {
		res, err := http.Get(s.URL + "/bar")
		if i == 4 {
			require.Equal(t, 500, res.StatusCode)
			break
		}

		require.Equal(t, err, nil)
		require.Equal(t, res.StatusCode, 200)
		body, _ := io.ReadAll(res.Body)
		require.Equal(t, string(body)[:13], `{"foo":"bar"}`)
	}
}

// TODO FIX THIS
// func TestUnmatched(t *testing.T) {
// 	defer after()
//
// 	// clear out any unmatchedRequests from other tests
// 	unmatchedRequests = []*http.Request{}
//
// 	// Intercept()
//
// 	_, err := http.Get(s.URL+"/unmatched")
// 	st.Reject(t, err, nil)
//
// 	unmatched := GetUnmatchedRequests()
// 	require.Equal(t, len(unmatched), 1)
// 	require.Equal(t, unmatched[0].URL.Host, "server.com")
// 	require.Equal(t, unmatched[0].URL.Path, "/unmatched")
// 	require.Equal(t, HasUnmatchedRequest(), true)
// }

// TODO FIX THIX
// func TestMultipleMocks(t *testing.T) {
// 	defer Disable()
//
// 	New(s.URL).
// 		Get("/foo").
// 		Reply(200).
// 		JSON(map[string]string{"value": "foo"})
//
// 	New(s.URL).
// 		Get("/bar").
// 		Reply(200).
// 		JSON(map[string]string{"value": "bar"})
//
// 	New(s.URL).
// 		Get("/baz").
// 		Reply(200).
// 		JSON(map[string]string{"value": "baz"})
//
// 	tests := []struct {
// 		path string
// 	}{
// 		{"/foo"},
// 		{"/bar"},
// 		{"/baz"},
// 	}
//
// 	for _, test := range tests {
// 		res, err := http.Get(s.URL + test.path)
// 		require.Equal(t, err, nil)
// 		require.Equal(t, res.StatusCode, 200)
// 		body, _ := io.ReadAll(res.Body)
// 		require.Equal(t, string(body)[:15], `{"value":"`+test.path[1:]+`"}`)
// 	}
//
// 	_, err := http.Get(s.URL+"/foo")
// 	st.Reject(t, err, nil)
// }

// func TestInterceptClient(t *testing.T) {
// 	defer after()
//
// 	New(s.URL).Reply(204)
// 	require.Equal(t, len(GetAll()), 1)
//
// 	req, err := http.NewRequest("GET", s.URL, nil)
// 	client := &http.Client{Transport: &http.Transport{}}
// 	InterceptClient(client)
//
// 	res, err := client.Do(req)
// 	require.Equal(t, err, nil)
// 	require.Equal(t, res.StatusCode, 204)
// }

// func TestRestoreClient(t *testing.T) {
// 	defer after()
//
// 	New(s.URL).Reply(204)
// 	require.Equal(t, len(GetAll()), 1)
//
// 	req, err := http.NewRequest("GET", s.URL, nil)
// 	client := &http.Client{Transport: &http.Transport{}}
// 	InterceptClient(client)
// 	trans := client.Transport
//
// 	res, err := client.Do(req)
// 	require.Equal(t, err, nil)
// 	require.Equal(t, res.StatusCode, 204)
//
// 	RestoreClient(client)
// 	st.Reject(t, trans, client.Transport)
// }

func TestMockRegExpMatching(t *testing.T) {
	t.Parallel()

	defer after()
	s := Server(t)
	New(s.URL).
		Post("/bar").
		MatchHeader("Authorization", "Bearer (.*)").
		BodyString(`{"foo":".*"}`).
		Reply(200).
		SetHeader("Server", "gock").
		JSON(map[string]string{"foo": "bar"})

	req, _ := http.NewRequest("POST", s.URL+"/bar", bytes.NewBuffer([]byte(`{"foo":"baz"}`)))
	req.Header.Set("Authorization", "Bearer s3cr3t")

	res, err := http.DefaultClient.Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	require.Equal(t, res.Header.Get("Server"), "gock")

	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body)[:13], `{"foo":"bar"}`)
}

//
// func TestObserve(t *testing.T) {
// 	// t.Parallel()
//
// 	defer after()
// 	var observedRequest *http.Request
// 	var observedMock Mock
// 	s := Server(t)
// 	Observe(func(request *http.Request, mock Mock) {
// 		observedRequest = request
// 		observedMock = mock
// 	})
// 	New(s.URL).Reply(200)
// 	req, _ := http.NewRequest("POST", s.URL, nil)
//
// 	http.DefaultClient.Do(req)
//
// 	require.NotNil(t, observedRequest)
// 	require.Contains(t, s.URL, observedRequest.Host)
// 	require.Contains(t, s.URL, observedMock.Request().URLStruct.Host)
// }

//
// func TestTryCreatingRacesInNew(s.URL) {
// 	defer after()
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			New(s.URL)
// 		}()
// 	}
// }

func after() {
	// Flush()
	// Disable()
}
