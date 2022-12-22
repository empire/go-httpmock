package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMatchHeaders(t *testing.T) {
	defer gock.Disable()

	gock.New("http://foo.com").
		MatchHeader("Authorization", "^foo bar$").
		MatchHeader("API", "1.[0-9]+").
		HeaderPresent("Accept").
		Reply(200).
		BodyString("foo foo")

	req, err := http.NewRequest("GET", "http://foo.com", nil)
	req.Header.Set("Authorization", "foo bar")
	req.Header.Set("API", "1.0")
	req.Header.Set("Accept", "text/plain")

	res, err := (&http.Client{}).Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}
