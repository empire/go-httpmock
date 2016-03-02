package gock

import (
	"bytes"
	"fmt"
	"github.com/nbio/st"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMockSimple(t *testing.T) {
	defer after()
	New("http://foo.com").Reply(201).JSON(map[string]string{"foo": "bar"})
	res, err := http.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}

func TestMockBodyStringResponse(t *testing.T) {
	defer after()
	New("http://foo.com").Reply(200).BodyString("foo bar")
	res, err := http.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo bar")
}

func TestMockBodyMatch(t *testing.T) {
	defer after()
	New("http://foo.com").BodyString("foo bar").Reply(201).BodyString("foo foo")
	res, err := http.Post("http://foo.com", "text/plain", bytes.NewBuffer([]byte("foo bar")))
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}

func TestMockBodyMatchJSON(t *testing.T) {
	defer after()
	New("http://foo.com").
		Post("/bar").
		JSON(map[string]string{"foo": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	res, err := http.Post("http://foo.com/bar", "application/json", bytes.NewBuffer([]byte(`{"foo":"bar"}`)))
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"bar":"foo"}`)
}

func TestMockMatchHeaders(t *testing.T) {
	defer after()
	New("http://foo.com").
		MatchHeader("Content-Type", "(.*)/plain").
		Reply(200).
		BodyString("foo foo")

	res, err := http.Post("http://foo.com", "text/plain", bytes.NewBuffer([]byte("foo bar")))
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}

func TestMockMap(t *testing.T) {
	defer after()

	mock := New("http://bar.com")
	mock.Map(func(req *http.Request) *http.Request {
		req.URL.Host = "bar.com"
		return req
	})
	mock.Reply(201).JSON(map[string]string{"foo": "bar"})

	res, err := http.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}

func TestMockFilter(t *testing.T) {
	defer after()

	mock := New("http://foo.com")
	mock.Filter(func(req *http.Request) bool {
		return req.URL.Host == "foo.com"
	})
	mock.Reply(201).JSON(map[string]string{"foo": "bar"})

	res, err := http.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}

func TestMockCounterDisabled(t *testing.T) {
	defer after()
	New("http://foo.com").Reply(204)
	st.Expect(t, len(GetAll()), 1)
	res, err := http.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 204)
	st.Expect(t, len(GetAll()), 0)
}

func TestMockEnableNetwork(t *testing.T) {
	defer after()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world")
	}))
	defer ts.Close()

	EnableNetworking()
	defer DisableNetworking()

	New(ts.URL).Reply(204)
	st.Expect(t, len(GetAll()), 1)

	res, err := http.Get(ts.URL)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 204)
	st.Expect(t, len(GetAll()), 0)

	res, err = http.Get(ts.URL)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
}

func TestMockEnableNetworkFilter(t *testing.T) {
	defer after()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world")
	}))
	defer ts.Close()

	EnableNetworking()
	defer DisableNetworking()

	NetworkingFilter(func(req *http.Request) bool {
		return strings.Contains(req.URL.Host, "127.0.0.1")
	})
	defer DisableNetworkingFilters()

	New(ts.URL).Reply(0).SetHeader("foo", "bar")
	st.Expect(t, len(GetAll()), 1)

	res, err := http.Get(ts.URL)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	st.Expect(t, res.Header.Get("foo"), "bar")
	st.Expect(t, len(GetAll()), 0)
}

func TestInterceptClient(t *testing.T) {
	defer after()

	New("http://foo.com").Reply(204)
	st.Expect(t, len(GetAll()), 1)

	req, err := http.NewRequest("GET", "http://foo.com", nil)
	client := &http.Client{Transport: &http.Transport{}}
	InterceptClient(client)

	res, err := client.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 204)
}

func after() {
	Flush()
	Disable()
}
