package test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

func TestReplyError(t *testing.T) {
	defer httpmock.Off()

	httpmock.New("http://foo.com").
		Get("/bar").
		ReplyError(errors.New("Error dude!"))

	_, err := http.Get("http://foo.com/bar")
	require.Equal(t, strings.HasSuffix(err.Error(), ": Error dude!"), true)

	// Verify that we don't have pending mocks
	require.Equal(t, httpmock.IsDone(), true)
}
