# httpmock [![Build Status](https://travis-ci.org/empire/go-httpmock.svg?branch=master)](https://travis-ci.org/empire/go-httpmock) [![GitHub release](https://img.shields.io/badge/version-v1.0-orange.svg?style=flat)](https://github.com/empire/go-httpmock/releases) [![GoDoc](https://godoc.org/github.com/empire/go-httpmock?status.svg)](https://godoc.org/github.com/empire/go-httpmock) [![Coverage Status](https://coveralls.io/repos/github/empire/go-httpmock/badge.svg?branch=master)](https://coveralls.io/github/empire/go-httpmock?branch=master) [![Go Report Card](https://img.shields.io/badge/go_report-A+-brightgreen.svg)](https://goreportcard.com/report/github.com/empire/go-httpmock) [![license](https://img.shields.io/badge/license-MIT-blue.svg)]()

Versatile HTTP mocking made easy in [Go](https://golang.org) that works with any `net/http` based stdlib implementation.

Heavily inspired by [nock](https://github.com/node-nock/nock).
There is also its Python port, [pook](https://github.com/h2non/pook).

To get started, take a look to the [examples](#examples).

## Features

- Simple, expressive, fluent API.
- Semantic API DSL for declarative HTTP mock declarations.
- Built-in helpers for easy JSON/XML mocking.
- Supports persistent and volatile TTL-limited mocks.
- Full regular expressions capable HTTP request mock matching.
- Designed for both testing and runtime scenarios.
- Match request by method, URL params, headers and bodies.
- Extensible and pluggable HTTP matching rules.
- Ability to switch between mock and real networking modes.
- Ability to filter/map HTTP requests for accurate mock matching.
- Supports map and filters to handle mocks easily.
- Wide compatible HTTP interceptor using `http.RoundTripper` interface.
- Works with any `net/http` compatible client, such as [gentleman](https://github.com/h2non/gentleman).
- Network timeout/cancelation delay simulation.
- Extensible and hackable API.
- Dependency free.

## Installation

```bash
go get -u github.com/empire/go-httpmock
```

## API

See [godoc reference](https://godoc.org/github.com/empire/go-httpmock) for detailed API documentation.

## How it mocks

1. Intercepts any HTTP outgoing request via `http.DefaultTransport` or custom `http.Transport` used by any `http.Client`.
2. Matches outgoing HTTP requests against a pool of defined HTTP mock expectations in FIFO declaration order.
3. If at least one mock matches, it will be used in order to compose the mock HTTP response.
4. If no mock can be matched, it will resolve the request with an error, unless real networking mode is enable, in which case a real HTTP request will be performed.

## Tips

#### Testing

Declare your mocks before you start declaring the concrete test logic:

```go
func TestFoo(t *testing.T) {
  defer httpmock.Off() // Flush pending mocks after test execution

  httpmock.New("http://server.com").
    Get("/bar").
    Reply(200).
    JSON(map[string]string{"foo": "bar"})

  // Your test code starts here...
}
```

#### Race conditions

If you're running concurrent code, be aware that your mocks are declared first to avoid unexpected
race conditions while configuring `httpmock` or intercepting custom HTTP clients.

`httpmock` is not fully thread-safe, but sensible parts are.
Any help making `httpmock` more reliable in this sense is appreciated.

#### Define complex mocks first

If you're mocking a bunch of mocks in the same test suite, it's recommended to define the more
concrete mocks first, and then the generic ones.

This approach usually avoids matching unexpected generic mocks (e.g: specific header, body payload...) instead of the generic ones that performs less complex matches.

#### Disable `httpmock` traffic interception once done

In other to minimize potential side effects within your test code, it's a good practice
disabling `httpmock` once you are done with your HTTP testing logic.

A Go idiomatic approach for doing this can be using it in a `defer` statement, such as:

```go
func TestGock (t *testing.T) {
	defer httpmock.Off()

	// ... my test code goes here
}
```

#### Intercept an `http.Client` just once

You don't need to intercept multiple times the same `http.Client` instance.

Just call `httpmock.InterceptClient(client)` once, typically at the beginning of your test scenarios.

#### Restore an `http.Client` after interception

**NOTE**: this is not required is you are using `http.DefaultClient` or `http.DefaultTransport`.

As a good testing pattern, you should call `httpmock.RestoreClient(client)` after running your test scenario, typically as after clean up hook.

You can also use a `defer` statement for doing it, as you do with `httpmock.Off()`, such as:

```go
func TestGock (t *testing.T) {
	defer httpmock.Off()
	defer httpmock.RestoreClient(client)

	// ... my test code goes here
}
```

## Examples

See [examples](https://github.com/empire/go-httpmock/tree/master/_examples) directory for more featured use cases.

#### Simple mocking via tests

```go
package test

import (
  "github.com/stretchr/testify/require"
  "github.com/empire/go-httpmock"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestSimple(t *testing.T) {
  defer httpmock.Off()

  httpmock.New("http://foo.com").
    Get("/bar").
    Reply(200).
    JSON(map[string]string{"foo": "bar"})

  res, err := http.Get("http://foo.com/bar")
  require.Equal(t, err, nil)
  require.Equal(t, res.StatusCode, 200)

  body, _ := ioutil.ReadAll(res.Body)
  require.Equal(t, string(body)[:13], `{"foo":"bar"}`)

  // Verify that we don't have pending mocks
  require.Equal(t, httpmock.IsDone(), true)
}
```

#### Request headers matching

```go
package test

import (
  "github.com/stretchr/testify/require"
  "github.com/empire/go-httpmock"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestMatchHeaders(t *testing.T) {
  defer httpmock.Off()

  httpmock.New("http://foo.com").
    MatchHeader("Authorization", "^foo bar$").
    MatchHeader("API", "1.[0-9]+").
    HeaderPresent("Accept").
    Reply(200).
    BodyString("foo foo")

  req, err := http.NewRequest("GET", "http://foo.com", nil)
  req.Header.Set("Authorization", "foo bar")
  req.Header.Set("API", "1.0")
  req.Header.Set("Accept", "text/plain")

  res, err := (&http.Client{}).Do(req)
  require.Equal(t, err, nil)
  require.Equal(t, res.StatusCode, 200)
  body, _ := ioutil.ReadAll(res.Body)
  require.Equal(t, string(body), "foo foo")

  // Verify that we don't have pending mocks
  require.Equal(t, httpmock.IsDone(), true)
}
```

#### Request param matching

```go
package test

import (
  "github.com/stretchr/testify/require"
  "github.com/empire/go-httpmock"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestMatchParams(t *testing.T) {
  defer httpmock.Off()

  httpmock.New("http://foo.com").
    MatchParam("page", "1").
    MatchParam("per_page", "10").
    Reply(200).
    BodyString("foo foo")

  req, err := http.NewRequest("GET", "http://foo.com?page=1&per_page=10", nil)

  res, err := (&http.Client{}).Do(req)
  require.Equal(t, err, nil)
  require.Equal(t, res.StatusCode, 200)
  body, _ := ioutil.ReadAll(res.Body)
  require.Equal(t, string(body), "foo foo")

  // Verify that we don't have pending mocks
  require.Equal(t, httpmock.IsDone(), true)
}
```

#### JSON body matching and response

```go
package test

import (
  "bytes"
  "github.com/stretchr/testify/require"
  "github.com/empire/go-httpmock"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestMockSimple(t *testing.T) {
  defer httpmock.Off()

  httpmock.New("http://foo.com").
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

  // Verify that we don't have pending mocks
  require.Equal(t, httpmock.IsDone(), true)
}
```

#### Mocking a custom http.Client and http.RoundTripper

```go
package test

import (
  "github.com/stretchr/testify/require"
  "github.com/empire/go-httpmock"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestClient(t *testing.T) {
  defer httpmock.Off()

  httpmock.New("http://foo.com").
    Reply(200).
    BodyString("foo foo")

  req, err := http.NewRequest("GET", "http://foo.com", nil)
  client := &http.Client{Transport: &http.Transport{}}
  httpmock.InterceptClient(client)

  res, err := client.Do(req)
  require.Equal(t, err, nil)
  require.Equal(t, res.StatusCode, 200)
  body, _ := ioutil.ReadAll(res.Body)
  require.Equal(t, string(body), "foo foo")

  // Verify that we don't have pending mocks
  require.Equal(t, httpmock.IsDone(), true)
}
```

#### Enable real networking

```go
package main

import (
  "fmt"
  "github.com/empire/go-httpmock"
  "io/ioutil"
  "net/http"
)

func main() {
  defer httpmock.Off()
  defer httpmock.DisableNetworking()

  httpmock.EnableNetworking()
  httpmock.New("http://httpbin.org").
    Get("/get").
    Reply(201).
    SetHeader("Server", "httpmock")

  res, err := http.Get("http://httpbin.org/get")
  if err != nil {
    fmt.Errorf("Error: %s", err)
  }

  // The response status comes from the mock
  fmt.Printf("Status: %d\n", res.StatusCode)
  // The server header comes from mock as well
  fmt.Printf("Server header: %s\n", res.Header.Get("Server"))
  // Response body is the original
  body, _ := ioutil.ReadAll(res.Body)
  fmt.Printf("Body: %s", string(body))
}
```

#### Debug intercepted http requests

```go
package main

import (
	"bytes"
	"github.com/empire/go-httpmock"
	"net/http"
)

func main() {
	defer httpmock.Off()
	httpmock.Observe(httpmock.DumpRequest)

	httpmock.New("http://foo.com").
		Post("/bar").
		MatchType("json").
		JSON(map[string]string{"foo": "bar"}).
		Reply(200)

	body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
	http.Post("http://foo.com/bar", "application/json", body)
}

```

## Hacking it!

You can easily hack `httpmock` defining custom matcher functions with own matching rules.

See [add matcher functions](https://github.com/empire/go-httpmock/blob/master/_examples/add_matchers/matchers.go) and [custom matching layer](https://github.com/empire/go-httpmock/blob/master/_examples/custom_matcher/matcher.go) examples for further details.

## License

MIT - Tomas Aparicio
