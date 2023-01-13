package test

import (
	"net/http"
	"testing"

	"github.com/empire/go-httpmock"
	"github.com/stretchr/testify/require"
)

// TODO use the httpmoc in main function
func Test_IsPending(t *testing.T) {
	t.Parallel()

	s := httpmock.Server(t)
	require := require.New(t)

	httpmock.New(s.URL).
		Get("/get").
		Reply(204).
		SetHeader("Server", "gock")

	require.Len(httpmock.Pending(t), 1, "pending mocks before request")
	require.True(httpmock.IsPending(t), "is pending before request")

	_, err := http.Get(s.URL + "/get")
	require.NoError(err)

	require.Len(httpmock.Pending(t), 0, "Pending mocks after request: %d\n")
	require.False(httpmock.IsPending(t))
}
