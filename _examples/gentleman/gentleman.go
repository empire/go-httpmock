package main

import (
	"fmt"

	"github.com/empire/go-httpmock"
	"gopkg.in/h2non/gentleman.v1"
	"gopkg.in/h2non/gentleman.v1/context"
)

// Usege example with gentleman HTTP client toolkit.
// See also: https://github.com/h2non/gentleman-mock
func main() {
	defer httpmock.Off()

	httpmock.New("http://httpbin.org").
		Get("/*").
		Reply(204).
		SetHeader("Server", "gock")

	cli := gentleman.New()

	cli.UseHandler("before dial", func(ctx *context.Context, h context.Handler) {
		httpmock.InterceptClient(ctx.Client)
		h.Next(ctx)
	})

	res, err := cli.Request().URL("http://httpbin.org/get").Send()
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}

	fmt.Printf("Status: %d\n", res.StatusCode)
	fmt.Printf("Server header: %s\n", res.Header.Get("Server"))
}
