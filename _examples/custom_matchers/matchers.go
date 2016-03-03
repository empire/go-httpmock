package main

import (
	"fmt"
	"gopkg.in/h2non/gock.v0"
	"net/http"
)

func main() {
	defer gock.Disable()

	gock.New("http://httpbin.org").
		Get("/").
		AddMatcher(func(req *http.Request, ereq *gock.Request) bool {
			return req.URL.Scheme == "http"
		}).
		AddMatcher(func(req *http.Request, ereq *gock.Request) bool {
			return req.Method == ereq.Method
		}).
		Reply(204).
		SetHeader("Server", "gock")

	res, err := http.Get("http://httpbin.org/get")
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}

	fmt.Printf("Status: %d\n", res.StatusCode)
	fmt.Printf("Server header: %s\n", res.Header.Get("Server"))
}
