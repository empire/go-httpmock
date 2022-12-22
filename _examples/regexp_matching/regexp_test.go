package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestRegExpMatching(t *testing.T) {
	defer httpmock.Disable()
	httpmock.New("http://foo.com").
		Post("/bar").
		MatchHeader("Authorization", "Bearer (.*)").
		BodyString(`{"foo":".*"}`).
		Reply(200).
		SetHeader("Server", "gock").
		JSON(map[string]string{"foo": "bar"})

	req, _ := http.NewRequest("POST", "http://foo.com/bar", bytes.NewBuffer([]byte(`{"foo":"baz"}`)))
	req.Header.Set("Authorization", "Bearer s3cr3t")

	res, err := http.DefaultClient.Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	require.Equal(t, res.Header.Get("Server"), "gock")
	body, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(body)[:13], `{"foo":"bar"}`)
}
