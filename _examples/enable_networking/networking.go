package test

import (
	"fmt"
	"gopkg.in/h2non/gock.v0"
	"net/http"
)

func main() {
	defer gock.Disable()
	defer gock.DisableNetworking()

	gock.EnableNetworking()
	gock.New("http://httpbin.org").
		Get("/get").
		Reply(201)

	res, err := http.Get("http://httpbin.org/get")
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}
	fmt.Printf("Status: %d\n", res.StatusCode)
}
