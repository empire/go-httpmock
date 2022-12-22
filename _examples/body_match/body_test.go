package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestMockSimple(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Post("/bar").
		MatchType("json").
		JSON(map[string]string{"foo": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
	res, err := http.Post("http://foo.com/bar", "application/json", body)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 201)

	resBody, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(resBody)[:13], `{"bar":"foo"}`)
}
