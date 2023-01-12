package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMatchURL(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)

	// TODO can we change the api to support this?
	// httpmock.New("http://(.*).com").

	httpmock.New(s.URL).
		Reply(200).
		BodyString("foo foo")

	res, err := http.Get(s.URL)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}
