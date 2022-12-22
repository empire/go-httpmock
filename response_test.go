package httpmock

import (
	"bytes"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewResponse(t *testing.T) {
	res := NewResponse()

	res.Status(200)
	require.Equal(t, res.StatusCode, 200)

	res.SetHeader("foo", "bar")
	require.Equal(t, res.Header.Get("foo"), "bar")

	res.Delay(1000 * time.Millisecond)
	require.Equal(t, res.ResponseDelay, 1000*time.Millisecond)

	// res.EnableNetworking()
	// require.Equal(t, res.UseNetwork, true)
}

func TestResponseStatus(t *testing.T) {
	res := NewResponse()
	require.Equal(t, res.StatusCode, 0)
	res.Status(200)
	require.Equal(t, res.StatusCode, 200)
}

func TestResponseType(t *testing.T) {
	res := NewResponse()
	res.Type("json")
	require.Equal(t, res.Header.Get("Content-Type"), "application/json")

	res = NewResponse()
	res.Type("xml")
	require.Equal(t, res.Header.Get("Content-Type"), "application/xml")

	res = NewResponse()
	res.Type("foo/bar")
	require.Equal(t, res.Header.Get("Content-Type"), "foo/bar")
}

func TestResponseSetHeader(t *testing.T) {
	res := NewResponse()
	res.SetHeader("foo", "bar")
	res.SetHeader("bar", "baz")
	require.Equal(t, res.Header.Get("foo"), "bar")
	require.Equal(t, res.Header.Get("bar"), "baz")
}

func TestResponseAddHeader(t *testing.T) {
	res := NewResponse()
	res.AddHeader("foo", "bar")
	res.AddHeader("foo", "baz")
	require.Equal(t, res.Header.Get("foo"), "bar")
	require.Equal(t, res.Header["Foo"][1], "baz")
}

func TestResponseSetHeaders(t *testing.T) {
	res := NewResponse()
	res.SetHeaders(map[string]string{"foo": "bar", "bar": "baz"})
	require.Equal(t, res.Header.Get("foo"), "bar")
	require.Equal(t, res.Header.Get("bar"), "baz")
}

func TestResponseBody(t *testing.T) {
	res := NewResponse()
	res.Body(bytes.NewBuffer([]byte("foo bar")))
	require.Equal(t, string(res.BodyBuffer), "foo bar")
}

func TestResponseBodyString(t *testing.T) {
	res := NewResponse()
	res.BodyString("foo bar")
	require.Equal(t, string(res.BodyBuffer), "foo bar")
}

func TestResponseFile(t *testing.T) {
	res := NewResponse()
	res.File("version.go")
	require.Equal(t, "package httpmock", string(res.BodyBuffer)[:16])
}

func TestResponseJSON(t *testing.T) {
	res := NewResponse()
	res.JSON(map[string]string{"foo": "bar"})
	require.Equal(t, string(res.BodyBuffer)[:13], `{"foo":"bar"}`)
	require.Equal(t, res.Header.Get("Content-Type"), "application/json")
}

func TestResponseXML(t *testing.T) {
	res := NewResponse()
	type xml struct {
		Data string `xml:"data"`
	}
	res.XML(xml{Data: "foo"})
	require.Equal(t, string(res.BodyBuffer), `<xml><data>foo</data></xml>`)
	require.Equal(t, res.Header.Get("Content-Type"), "application/xml")
}

func TestResponseMap(t *testing.T) {
	res := NewResponse()
	require.Equal(t, len(res.Mappers), 0)
	res.Map(func(res *http.Response) *http.Response {
		return res
	})
	require.Equal(t, len(res.Mappers), 1)
}

func TestResponseFilter(t *testing.T) {
	res := NewResponse()
	require.Equal(t, len(res.Filters), 0)
	res.Filter(func(res *http.Response) bool {
		return true
	})
	require.Equal(t, len(res.Filters), 1)
}

func TestResponseSetError(t *testing.T) {
	res := NewResponse()
	require.Equal(t, res.Error, nil)
	res.SetError(errors.New("foo error"))
	require.Equal(t, res.Error.Error(), "foo error")
}

func TestResponseDelay(t *testing.T) {
	res := NewResponse()
	require.Equal(t, res.ResponseDelay, 0*time.Microsecond)
	res.Delay(100 * time.Millisecond)
	require.Equal(t, res.ResponseDelay, 100*time.Millisecond)
}

//
// func TestResponseEnableNetworking(t *testing.T) {
// 	res := NewResponse()
// 	require.Equal(t, res.UseNetwork, false)
// 	res.EnableNetworking()
// 	require.Equal(t, res.UseNetwork, true)
// }

func TestResponseDone(t *testing.T) {
	res := NewResponse()
	res.Mock = &Mocker{request: &Request{Counter: 1}, disabler: new(disabler)}
	require.Equal(t, res.Done(), false)
	res.Mock.Disable()
	require.Equal(t, res.Done(), true)
}
