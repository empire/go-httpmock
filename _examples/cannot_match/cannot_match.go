package main

import (
	"fmt"
	"net/http"

	"github.com/empire/go-httpmock"
)

func main() {
	// gock enabled but cannot match any mock
	httpmock.New("http://httpbin.org").
		Get("/foo").
		Reply(201).
		SetHeader("Server", "gock")

	_, err := http.Get("http://httpbin.org/bar")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
