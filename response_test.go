package gock

import (
	"github.com/nbio/st"
	"testing"
	"time"
)

func TestNewResponse(t *testing.T) {
	res := NewResponse()

	res.Status(200)
	st.Expect(t, res.StatusCode, 200)

	res.SetHeader("foo", "bar")
	st.Expect(t, res.Header.Get("foo"), "bar")

	res.Delay(1000 * time.Millisecond)
	st.Expect(t, res.ResponseDelay, 1000*time.Millisecond)

	res.EnableNetworking()
	st.Expect(t, res.UseNetwork, true)
}
