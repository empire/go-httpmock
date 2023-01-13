package test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestReplyError(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)

	httpmock.New(s.URL).
		Get("/bar").
		ReplyError(errors.New("Error dude!"))

	resp, err := http.Get(s.URL + "/bar")
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotImplemented, resp.StatusCode)
	require.Equal(t, "Error dude!", string(body))

	// Verify that we don't have pending mocks
	require.True(t, httpmock.IsDone(t))
}
