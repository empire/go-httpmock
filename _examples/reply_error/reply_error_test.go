package test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestReplyError(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Get("/bar").
		ReplyError(errors.New("Error dude!"))

	_, err := http.Get("http://foo.com/bar")
	st.Expect(t, strings.HasSuffix(err.Error(), ": Error dude!"), true)

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}
