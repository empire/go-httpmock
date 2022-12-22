package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMatchURL(t *testing.T) {
	defer gock.Disable()

	gock.New("http://(.*).com").
		Reply(200).
		BodyString("foo foo")

	res, err := http.Get("http://foo.com")
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}
