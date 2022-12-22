package main

import (
	"bytes"
	"net/http"

	"github.com/empire/go-httpmock"
)

func main() {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	gock.New("http://foo.com").
		Post("/bar").
		MatchType("json").
		JSON(map[string]string{"foo": "bar"}).
		Reply(200)

	body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
	http.Post("http://foo.com/bar", "application/json", body)
}
