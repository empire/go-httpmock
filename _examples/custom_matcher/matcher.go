package main

import (
	"fmt"
	"net/http"

	"github.com/empire/go-httpmock"
)

func main() {
	defer httpmock.Off()

	// Create a new custom matcher with HTTP headers only matchers
	matcher := httpmock.NewBasicMatcher()

	// Add a custom match function
	matcher.Add(func(req *http.Request, ereq *httpmock.Request) (bool, error) {
		return req.URL.Scheme == "http", nil
	})

	// Define the mock
	httpmock.New("http://httpbin.org").
		SetMatcher(matcher).
		Get("/").
		Reply(204).
		SetHeader("Server", "gock")

	res, err := http.Get("http://httpbin.org/get")
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}

	fmt.Printf("Status: %d\n", res.StatusCode)
	fmt.Printf("Server header: %s\n", res.Header.Get("Server"))
}
