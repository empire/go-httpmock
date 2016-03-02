package gock

import "net/http"

// Mock represents the required interface that must
// be implemented by HTTP mock instances.
type Mock interface {
	// Disable disables the current mock manually.
	Disable()

	// Done returns true if the current mock is disabled.
	Done() bool

	// Request returns the mock Request instance.
	Request() *Request

	// Response returns the mock Response instance.
	Response() *Response

	// Match matches the given http.Request with the current mock.
	Match(*http.Request) (bool, error)

	// AddMatcher adds a new matcher function.
	AddMatcher(MatchFunc)
}

// Mocker implements a Mock capable interface providing
// a default mock configuration used internally to store mocks.
type Mocker struct {
	// disabled stores if the current mock is disabled.
	disabled bool

	// matcher stores a Matcher capable instance to match the given http.Request.
	matcher Matcher

	// request stores the mock Request to match.
	request *Request

	// response stores the mock Response to use in case of match.
	response *Response
}

// NewMock creates a new HTTP mock based on the given request and response instances.
// It's mostly used internally.
func NewMock(req *Request, res *Response) *Mocker {
	mock := &Mocker{
		request:  req,
		response: res,
		matcher:  DefaultMatcher,
	}
	res.Mock = mock
	req.Mock = mock
	req.Response = res
	return mock
}

// Disable disables the current mock manually.
func (e *Mocker) Disable() {
	e.disabled = true
}

// Done returns true in case that the current mock
// instance is disabled and therefore must be removed.
func (e *Mocker) Done() bool {
	return e.disabled || (!e.request.Persisted && e.request.Counter == 0)
}

// Request returns the Request instance
// configured for the current HTTP mock.
func (e *Mocker) Request() *Request {
	return e.request
}

// Response returns the Response instance
// configured for the current HTTP mock.
func (e *Mocker) Response() *Response {
	return e.response
}

// Match matches the given http.Request with the current Request
// mock expectation, returning true if matches.
func (e *Mocker) Match(req *http.Request) (bool, error) {
	if e.disabled {
		return false, nil
	}

	// Map
	for _, mapper := range e.request.Mappers {
		if treq := mapper(req); treq != nil {
			req = treq
		}
	}

	// Filter
	for _, filter := range e.request.Filters {
		if !filter(req) {
			return false, nil
		}
	}

	// Match
	matches, err := e.matcher.Match(req, e.request)
	if matches {
		e.decrement()
	}

	return matches, err
}

// SetMatcher sets a new matcher implementation
// for the current mock expectation.
func (e *Mocker) SetMatcher(matcher Matcher) {
	e.matcher = matcher
}

// AddMatcher adds a new matcher function
// for the current mock expectation.
func (e *Mocker) AddMatcher(fn MatchFunc) {
	e.matcher.Add(fn)
}

// decrement decrements the current mock Request counter.
func (e *Mocker) decrement() {
	if e.request.Persisted {
		return
	}

	e.request.Counter--
	if e.request.Counter == 0 {
		e.disabled = true
	}
}
