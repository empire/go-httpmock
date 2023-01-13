package test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMockSimple(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)

	httpmock.New(s.URL).
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
	req, err := http.NewRequest("POST", s.URL+"/bar", &compressed)
	require.Equal(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"bar":"foo"}`, string(resBody))
}
