# gock [![Build Status](https://travis-ci.org/h2non/gock.png)](https://travis-ci.org/h2non/gock) [![GitHub release](https://img.shields.io/badge/version-0.1.0-orange.svg?style=flat)](https://github.com/h2non/gock/releases) [![GoDoc](https://godoc.org/github.com/h2non/gock?status.svg)](https://godoc.org/github.com/h2non/gock) [![Coverage Status](https://coveralls.io/repos/github/h2non/gock/badge.svg?branch=master)](https://coveralls.io/github/h2non/gock?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/h2non/gock)](https://goreportcard.com/report/github.com/h2non/gock)

Versatile HTTP mocking made simple for [Go](https://golang.org). 
Heavily inspired by [nock](https://github.com/pgte/nock).

Take a look to the [examples](#examples) to get started.

Note: still beta, needs more docs and test coverage.

## Features

- Simple, expressive, fluent API.
- Semantic DSL for easy HTTP mocks definition.
- Built-in helpers for easy JSON/XML mocking.
- Supports persistent and volatile mocks.
- Full regexp capable HTTP request matching.
- Designed for both testing and runtime scenarios.
- Match request by method, URL params, headers and bodies.
- Extensible HTTP matching rules.
- Supports map and filters to wotk with mocks.
- Ability to switch between mock and real networking modes.
- Unobstructure HTTP interceptor based on `http.RoundTripper`.
- Network delay simulation (beta).
- Extensible and hackable API.

## Installation

```bash
go get -u gopkg.in/h2non/gock.v0
```

## API

See [godoc reference](https://godoc.org/github.com/h2non/gock) for detailed API documentation.

## Examples

#### Simple mocking via tests

```go
package test

import (
  "github.com/nbio/st"
  "gopkg.in/h2non/gock.v0"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestSimple(t *testing.T) {
  defer gock.Disable()
  gock.New("http://foo.com").
    Get("/bar").
    Reply(200).
    JSON(map[string]string{"foo": "bar"})

  res, err := http.Get("http://foo.com/bar")
  st.Expect(t, err, nil)
  st.Expect(t, res.StatusCode, 200)

  body, _ := ioutil.ReadAll(res.Body)
  st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}
```

#### Request headers matching

```go
package test

import (
  "github.com/nbio/st"
  "gopkg.in/h2non/gock.v0"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestMatchHeaders(t *testing.T) {
  defer gock.Disable()

  gock.New("http://foo.com").
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
  st.Expect(t, err, nil)
  st.Expect(t, res.StatusCode, 200)
  body, _ := ioutil.ReadAll(res.Body)
  st.Expect(t, string(body), "foo foo")
}
```

#### JSON body matching and response

```go
package test

import (
  "bytes"
  "github.com/nbio/st"
  "gopkg.in/h2non/gock.v0"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestMockSimple(t *testing.T) {
  defer gock.Disable()
  gock.New("http://foo.com").
    Post("/bar").
    JSON(map[string]string{"foo": "bar"}).
    Reply(201).
    JSON(map[string]string{"bar": "foo"})

  body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
  res, err := http.Post("http://foo.com/bar", "application/json", body)
  st.Expect(t, err, nil)
  st.Expect(t, res.StatusCode, 201)

  resBody, _ := ioutil.ReadAll(res.Body)
  st.Expect(t, string(resBody)[:13], `{"bar":"foo"}`)
}
```

#### Mocking using a custom http.Client

```go
package test

import (
  "github.com/nbio/st"
  "gopkg.in/h2non/gock.v0"
  "io/ioutil"
  "net/http"
  "testing"
)

func TestClient(t *testing.T) {
  defer gock.Disable()

  gock.New("http://foo.com").
    Reply(200).
    BodyString("foo foo")

  req, err := http.NewRequest("GET", "http://foo.com", nil)
  client := &http.Client{Transport: &http.Transport{}}
  gock.InterceptClient(client)

  res, err := client.Do(req)
  st.Expect(t, err, nil)
  st.Expect(t, res.StatusCode, 200)
  body, _ := ioutil.ReadAll(res.Body)
  st.Expect(t, string(body), "foo foo")
}
```

#### Enable real networking

```go
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
```

## License 

MIT - Tomas Aparicio
