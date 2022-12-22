package main

import (
	"fmt"
	"github.com/empire/go-httpmock"
	"net/http"
)

func main() {
	defer gock.Off()

	gock.New("http://httpbin.org").
		Get("/get").
		Filter(func(req *http.Request) bool { return req.URL.Host == "httpbin.org" }).
		Filter(func(req *http.Request) bool { return req.URL.Path == "/get" }).
		Reply(204).
		SetHeader("Server", "gock")

	res, err := http.Get("http://httpbin.org/get")
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}

	fmt.Printf("Status: %d\n", res.StatusCode)
	fmt.Printf("Server header: %s\n", res.Header.Get("Server"))
}
