package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)

	httpmock.New(s.URL).
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	res, err := http.Get(s.URL + "/bar")
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)

	body, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(body)[:13], `{"foo":"bar"}`)

	// Verify that we don't have pending mocks
	require.True(t, httpmock.IsDone(t))
}
