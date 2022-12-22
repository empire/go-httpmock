package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Reply(200).
		BodyString("foo foo")

	req, err := http.NewRequest("GET", "http://foo.com", nil)
	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	res, err := client.Do(req)
	require.Equal(t, err, nil)
	require.Equal(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	require.Equal(t, string(body), "foo foo")
}
