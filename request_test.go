package httpmock

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest()
	req.URL("http://foo.com")
	require.Equal(t, req.URLStruct.Host, "foo.com")
	require.Equal(t, req.URLStruct.Scheme, "http")
	req.MatchHeader("foo", "bar")
	require.Equal(t, req.Header.Get("foo"), "bar")
}

func TestRequestSetURL(t *testing.T) {
	req := NewRequest()
	req.URL("http://foo.com")
	req.SetURL(&url.URL{Host: "bar.com", Path: "/foo"})
	require.Equal(t, req.URLStruct.Host, "bar.com")
	require.Equal(t, req.URLStruct.Path, "/foo")
}

func TestRequestPath(t *testing.T) {
	req := NewRequest()
	req.URL("http://foo.com")
	req.Path("/foo")
	require.Equal(t, req.URLStruct.Scheme, "http")
	require.Equal(t, req.URLStruct.Host, "foo.com")
	require.Equal(t, req.URLStruct.Path, "/foo")
}

func TestRequestBody(t *testing.T) {
	req := NewRequest()
	req.Body(bytes.NewBuffer([]byte("foo bar")))
	require.Equal(t, string(req.BodyBuffer), "foo bar")
}

func TestRequestBodyString(t *testing.T) {
	req := NewRequest()
	req.BodyString("foo bar")
	require.Equal(t, string(req.BodyBuffer), "foo bar")
}

func TestRequestFile(t *testing.T) {
	req := NewRequest()
	req.File("version.go")
	require.Equal(t, string(req.BodyBuffer)[:12], "package httpmock")
}

func TestRequestJSON(t *testing.T) {
	req := NewRequest()
	req.JSON(map[string]string{"foo": "bar"})
	require.Equal(t, string(req.BodyBuffer)[:13], `{"foo":"bar"}`)
	require.Equal(t, req.Header.Get("Content-Type"), "application/json")
}

func TestRequestXML(t *testing.T) {
	req := NewRequest()
	type xml struct {
		Data string `xml:"data"`
	}
	req.XML(xml{Data: "foo"})
	require.Equal(t, string(req.BodyBuffer), `<xml><data>foo</data></xml>`)
	require.Equal(t, req.Header.Get("Content-Type"), "application/xml")
}

func TestRequestMatchType(t *testing.T) {
	req := NewRequest()
	req.MatchType("json")
	require.Equal(t, req.Header.Get("Content-Type"), "application/json")

	req = NewRequest()
	req.MatchType("html")
	require.Equal(t, req.Header.Get("Content-Type"), "text/html")

	req = NewRequest()
	req.MatchType("foo/bar")
	require.Equal(t, req.Header.Get("Content-Type"), "foo/bar")
}

func TestRequestBasicAuth(t *testing.T) {
	req := NewRequest()
	req.BasicAuth("bob", "qwerty")
	require.Equal(t, req.Header.Get("Authorization"), "Basic Ym9iOnF3ZXJ0eQ==")
}

func TestRequestMatchHeader(t *testing.T) {
	req := NewRequest()
	req.MatchHeader("foo", "bar")
	req.MatchHeader("bar", "baz")
	req.MatchHeader("UPPERCASE", "bat")
	req.MatchHeader("Mixed-CASE", "foo")

	require.Equal(t, req.Header.Get("foo"), "bar")
	require.Equal(t, req.Header.Get("bar"), "baz")
	require.Equal(t, req.Header.Get("UPPERCASE"), "bat")
	require.Equal(t, req.Header.Get("Mixed-CASE"), "foo")
}

func TestRequestHeaderPresent(t *testing.T) {
	req := NewRequest()
	req.HeaderPresent("foo")
	req.HeaderPresent("bar")
	req.HeaderPresent("UPPERCASE")
	req.HeaderPresent("Mixed-CASE")
	require.Equal(t, req.Header.Get("foo"), ".*")
	require.Equal(t, req.Header.Get("bar"), ".*")
	require.Equal(t, req.Header.Get("UPPERCASE"), ".*")
	require.Equal(t, req.Header.Get("Mixed-CASE"), ".*")
}

func TestRequestMatchParam(t *testing.T) {
	req := NewRequest()
	req.MatchParam("foo", "bar")
	req.MatchParam("bar", "baz")
	require.Equal(t, req.URLStruct.Query().Get("foo"), "bar")
	require.Equal(t, req.URLStruct.Query().Get("bar"), "baz")
}

func TestRequestMatchParams(t *testing.T) {
	req := NewRequest()
	req.MatchParams(map[string]string{"foo": "bar", "bar": "baz"})
	require.Equal(t, req.URLStruct.Query().Get("foo"), "bar")
	require.Equal(t, req.URLStruct.Query().Get("bar"), "baz")
}

func TestRequestPresentParam(t *testing.T) {
	req := NewRequest()
	req.ParamPresent("key")
	require.Equal(t, req.URLStruct.Query().Get("key"), ".*")
}

func TestRequestPathParam(t *testing.T) {
	req := NewRequest()
	req.PathParam("key", "value")
	require.Equal(t, req.PathParams["key"], "value")
}

func TestRequestPersist(t *testing.T) {
	req := NewRequest()
	require.Equal(t, req.Persisted, false)
	req.Persist()
	require.Equal(t, req.Persisted, true)
}

func TestRequestTimes(t *testing.T) {
	req := NewRequest()
	require.Equal(t, req.Counter, 1)
	req.Times(3)
	require.Equal(t, req.Counter, 3)
}

func TestRequestMap(t *testing.T) {
	req := NewRequest()
	require.Equal(t, len(req.Mappers), 0)
	req.Map(func(req *http.Request) *http.Request {
		return req
	})
	require.Equal(t, len(req.Mappers), 1)
}

func TestRequestFilter(t *testing.T) {
	req := NewRequest()
	require.Equal(t, len(req.Filters), 0)
	req.Filter(func(req *http.Request) bool {
		return true
	})
	require.Equal(t, len(req.Filters), 1)
}

// func TestRequestEnableNetworking(t *testing.T) {
// 	req := NewRequest()
// 	req.Response = &Response{}
// 	require.Equal(t, req.Response.UseNetwork, false)
// 	req.EnableNetworking()
// 	require.Equal(t, req.Response.UseNetwork, true)
// }

func TestRequestResponse(t *testing.T) {
	req := NewRequest()
	res := NewResponse()
	req.Response = res
	chain := req.Reply(200)
	require.Equal(t, chain, res)
	require.Equal(t, chain.StatusCode, 200)
}

func TestRequestReplyFunc(t *testing.T) {
	req := NewRequest()
	res := NewResponse()
	req.Response = res
	chain := req.ReplyFunc(func(r *Response) {
		r.Status(204)
	})
	require.Equal(t, chain, res)
	require.Equal(t, chain.StatusCode, 204)
}

func TestRequestMethods(t *testing.T) {
	req := NewRequest()
	req.Get("/foo")
	require.Equal(t, req.Method, "GET")
	require.Equal(t, req.URLStruct.Path, "/foo")

	req = NewRequest()
	req.Post("/foo")
	require.Equal(t, req.Method, "POST")
	require.Equal(t, req.URLStruct.Path, "/foo")

	req = NewRequest()
	req.Put("/foo")
	require.Equal(t, req.Method, "PUT")
	require.Equal(t, req.URLStruct.Path, "/foo")

	req = NewRequest()
	req.Delete("/foo")
	require.Equal(t, req.Method, "DELETE")
	require.Equal(t, req.URLStruct.Path, "/foo")

	req = NewRequest()
	req.Patch("/foo")
	require.Equal(t, req.Method, "PATCH")
	require.Equal(t, req.URLStruct.Path, "/foo")

	req = NewRequest()
	req.Head("/foo")
	require.Equal(t, req.Method, "HEAD")
	require.Equal(t, req.URLStruct.Path, "/foo")
}

func TestRequestSetMatcher(t *testing.T) {
	defer after()

	matcher := NewEmptyMatcher()
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Host == "foo.com", nil
	})
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.Header.Get("foo") == "bar", nil
	})
	ereq := NewRequest()
	mock := NewMock(ereq, &Response{})
	mock.SetMatcher(matcher)
	ereq.Mock = mock

	headers := make(http.Header)
	headers.Set("foo", "bar")
	req := &http.Request{
		URL:    &url.URL{Host: "foo.com", Path: "/bar"},
		Header: headers,
	}

	match, err := ereq.Mock.Match(req)
	require.Equal(t, err, nil)
	require.Equal(t, match, true)
}

func TestRequestAddMatcher(t *testing.T) {
	defer after()

	ereq := NewRequest()
	mock := NewMock(ereq, &Response{})
	mock.matcher = NewMatcher()
	ereq.Mock = mock

	ereq.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Host == "foo.com", nil
	})
	ereq.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		return req.Header.Get("foo") == "bar", nil
	})

	headers := make(http.Header)
	headers.Set("foo", "bar")
	req := &http.Request{
		URL:    &url.URL{Host: "foo.com", Path: "/bar"},
		Header: headers,
	}

	match, err := ereq.Mock.Match(req)
	require.Equal(t, err, nil)
	require.Equal(t, match, true)
}
