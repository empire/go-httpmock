package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMatchQueryParams(t *testing.T) {
	defer gock.Disable()

	gock.New("http://foo.com").
		MatchParam("foo", "^bar$").
		MatchParam("bar", "baz").
		ParamPresent("baz").
		Reply(200).
		BodyString("foo foo")

	req, err := http.NewRequest("GET", "http://foo.com?foo=bar&bar=baz&baz=foo", nil)
	res, err := (&http.Client{}).Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}
