package test

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMockSimple(t *testing.T) {
	defer httpmock.Off()

	httpmock.New("http://foo.com").
		Post("/bar").
		MatchType("json").
		Compression("gzip").
		JSON(map[string]string{"foo": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write([]byte(`{"foo":"bar"}`))
	w.Close()
	req, err := http.NewRequest("POST", "http://foo.com/bar", &compressed)
	require.Equal(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)

	resBody, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(resBody)[:13], `{"bar":"foo"}`)
}
