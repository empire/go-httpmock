package gock

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Server(t *testing.T) {
	require := require.New(t)
	s := Server(t)

	New(s.URL).
		Get("/").
		Reply(200).
		JSON(map[string]any{
			"i": 2,
		})

	resp, body := SendRequestAndGetResponse(t, http.MethodGet, s, "/", nil, map[string]string{})

	require.Equal(200, resp.StatusCode)
	require.JSONEq(`{"i": 2}`, string(body))
}

func SendRequestAndGetResponse(t *testing.T, method string, server *httptest.Server, path string, body io.Reader, header map[string]string) (*http.Response, []byte) {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), method, server.URL+path, body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := server.Client()

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error making http call: %v", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response: %v", err)
	}

	return resp, b
}
