package gock

import (
	"github.com/nbio/st"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest()

	req.URL("http://foo.com")
	st.Expect(t, req.URLStruct.Host, "foo.com")
	st.Expect(t, req.URLStruct.Scheme, "http")

	req.MatchHeader("foo", "bar")
	st.Expect(t, req.Header.Get("foo"), "bar")
}
