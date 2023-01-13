package test

import (
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func Test_CannotMatch(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)

	httpmock.New(s.URL).
		Get("/foo").
		Reply(201).
		SetHeader("Server", "gock")

	resp, err := http.Get(s.URL + "/bar")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotImplemented, resp.StatusCode)
}
