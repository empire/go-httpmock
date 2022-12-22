package main

import (
	"fmt"
	"net/http"

	"github.com/empire/go-httpmock"
)

func main() {
	defer httpmock.Off()

	httpmock.New("http://httpbin.org").
		Get("/get").
		Reply(204).
		SetHeader("Server", "gock")

	fmt.Printf("Pending mocks before request: %d\n", len(httpmock.Pending()))
	fmt.Printf("Is pending before request: %#v\n", httpmock.IsPending())

	_, err := http.Get("http://httpbin.org/get")
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}

	fmt.Printf("Pending mocks after request: %d\n", len(httpmock.Pending()))
	fmt.Printf("Is pending: %#v\n", httpmock.IsPending())
}
