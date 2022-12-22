package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMultipleMocks(t *testing.T) {
	defer httpmock.Disable()

	httpmock.New("http://server.com").
		Get("/foo").
		Reply(200).
		JSON(map[string]string{"value": "foo"})

	httpmock.New("http://server.com").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"value": "bar"})

	httpmock.New("http://server.com").
		Get("/baz").
		Reply(200).
		JSON(map[string]string{"value": "baz"})

	tests := []struct {
		path string
	}{
		{"/foo"},
		{"/bar"},
		{"/baz"},
	}

	for _, test := range tests {
		res, err := http.Get("http://server.com" + test.path)
		require.Equal(t, err, nil)
		require.Equal(t, res.StatusCode, 200)
		body, _ := ioutil.ReadAll(res.Body)
		require.Equal(t, string(body)[:15], `{"value":"`+test.path[1:]+`"}`)
	}

	// Failed request after mocks expires
	_, err := http.Get("http://server.com/foo")
	st.Reject(t, err, nil)
}
