package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMultipleMocks(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)

	httpmock.New(s.URL).
		Get("/foo").
		Reply(200).
		JSON(map[string]string{"value": "foo"})

	httpmock.New(s.URL).
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"value": "bar"})

	httpmock.New(s.URL).
		Get("/baz").
		Reply(200).
		JSON(map[string]string{"value": "baz"})

	tests := []struct {
		path string
	}{
		{"/bar"},
		{"/foo"},
		{"/baz"},
	}

	for _, test := range tests {
		res, err := http.Get(s.URL + test.path)
		require.Equal(t, err, nil)
		require.Equal(t, res.StatusCode, 200)
		body, _ := io.ReadAll(res.Body)
		require.Equal(t, string(body)[:15], `{"value":"`+test.path[1:]+`"}`)
	}

	// Failed request after mocks expires
	resp, err := http.Get(s.URL + "/foo")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotImplemented, resp.StatusCode)
}
