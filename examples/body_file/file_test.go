package test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMockBodyFile(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)

	httpmock.New(s.URL).
		Post("/bar").
		MatchType("json").
		File("data.json").
		Reply(201).
		File("response.json")

	body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
	res, err := http.Post(s.URL+"/bar", "application/json", body)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)

	resBody, _ := io.ReadAll(res.Body)
	require.JSONEq(t, `{"bar":"foo"}`, string(resBody))
}
