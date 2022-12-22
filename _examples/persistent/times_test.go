package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestTimes(t *testing.T) {
	defer gock.Disable()
	gock.New("http://127.0.0.1:1234").
		Get("/bar").
		Times(4).
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	for i := 0; i < 5; i++ {
		res, err := http.Get("http://127.0.0.1:1234/bar")
		if i == 4 {
			st.Reject(t, err, nil)
			break
		}

		require.Equal(t, err, nil)
		require.Equal(t, res.StatusCode, 200)
		body, _ := ioutil.ReadAll(res.Body)
		require.Equal(t, string(body)[:13], `{"foo":"bar"}`)
	}
}
