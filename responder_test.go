package gock

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResponder(t *testing.T) {
	defer after()
	s := Server(t)
	mres := New(s.URL).Reply(200).BodyString("foo")
	req := &http.Request{}

	res, err := Responder(req, mres, nil)
	require.Equal(t, err, nil)
	require.Equal(t, res.Status, "200 OK")
	require.Equal(t, res.StatusCode, 200)

	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body), "foo")
}

func TestResponder_ReadTwice(t *testing.T) {
	defer after()
	s := Server(t)
	mres := New(s.URL).Reply(200).BodyString("foo")
	req := &http.Request{}

	res, err := Responder(req, mres, nil)
	require.Equal(t, err, nil)
	require.Equal(t, res.Status, "200 OK")
	require.Equal(t, res.StatusCode, 200)

	body, _ := io.ReadAll(res.Body)
	require.Equal(t, string(body), "foo")

	body, err = io.ReadAll(res.Body)
	require.Equal(t, err, nil)
	require.Equal(t, body, []byte{})
}

func TestResponderSupportsMultipleHeadersWithSameKey(t *testing.T) {
	defer after()
	s := Server(t)
	mres := New(s.URL).
		Reply(200).
		AddHeader("Set-Cookie", "a=1").
		AddHeader("Set-Cookie", "b=2")
	req := &http.Request{}

	res, err := Responder(req, mres, nil)
	require.Equal(t, err, nil)
	require.Equal(t, res.Header, http.Header{"Set-Cookie": []string{"a=1", "b=2"}})
}

func TestResponderError(t *testing.T) {
	defer after()
	s := Server(t)
	mres := New(s.URL).ReplyError(errors.New("error"))
	req := &http.Request{}

	res, err := Responder(req, mres, nil)
	require.Equal(t, err.Error(), "error")
	require.Equal(t, res == nil, true)
}

func TestResponderCancelledContext(t *testing.T) {
	defer after()
	s := Server(t)
	mres := New(s.URL).Get("").Reply(200).Delay(20 * time.Millisecond).BodyString("foo")

	// create a context and schedule a call to cancel in 10ms
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://foo.com", nil)

	res, err := Responder(req, mres, nil)

	// verify that we got a context cancellation error and nil response
	require.Equal(t, err, context.Canceled)
	require.Equal(t, res == nil, true)
}

func TestResponderExpiredContext(t *testing.T) {
	defer after()
	s := Server(t)
	mres := New(s.URL).Get("").Reply(200).Delay(20 * time.Millisecond).BodyString("foo")

	// create a context that is set to expire in 10ms
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://foo.com", nil)

	res, err := Responder(req, mres, nil)

	// verify that we got a context cancellation error and nil response
	require.Equal(t, err, context.DeadlineExceeded)
	require.Equal(t, res == nil, true)
}

func TestResponderPreExpiredContext(t *testing.T) {
	defer after()
	s := Server(t)
	mres := New(s.URL).Get("").Reply(200).BodyString("foo")

	// create a context and wait to ensure it is expired
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Microsecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://foo.com", nil)

	res, err := Responder(req, mres, nil)

	// verify that we got a context cancellation error and nil response
	require.Equal(t, err, context.DeadlineExceeded)
	require.Equal(t, res == nil, true)
}
