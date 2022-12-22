package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestPersistent(t *testing.T) {
	defer gock.Disable()
	gock.New("http://foo.com").
		Get("/bar").
		Persist().
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	for i := 0; i < 5; i++ {
		res, err := http.Get("http://foo.com/bar")
		require.Equal(t, err, nil)
		require.Equal(t, res.StatusCode, 200)
		body, _ := ioutil.ReadAll(res.Body)
		require.Equal(t, string(body)[:13], `{"foo":"bar"}`)
	}

	// Verify that we don't have pending mocks
	require.Equal(t, gock.IsDone(), true)
}
