package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMatchQueryParams(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)

	httpmock.New(s.URL).
		MatchParam("foo", "^bar$").
		MatchParam("bar", "baz").
		ParamPresent("baz").
		Reply(200).
		BodyString("foo foo")

	req, err := http.NewRequest("GET", s.URL+"?foo=bar&bar=baz&baz=foo", nil)
	require.NoError(t, err)

	res, err := (&http.Client{}).Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}
