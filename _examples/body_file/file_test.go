package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMockBodyFile(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Post("/bar").
		MatchType("json").
		File("data.json").
		Reply(201).
		File("response.json")

	body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
	res, err := http.Post("http://foo.com/bar", "application/json", body)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)

	resBody, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(resBody)[:13], `{"bar":"foo"}`)
}
