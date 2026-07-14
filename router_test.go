package router_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-amwk/core"
	"github.com/go-amwk/engine"
	"github.com/go-amwk/router"
)

// mockRequest implements core.Request with minimal functionality for tests.
type mockRequest struct {
	method string
	path   string
	vals   map[string]string
}

func (r *mockRequest) Body() (io.ReadCloser, error)        { return nil, nil }
func (r *mockRequest) ClientIP() string                    { return "" }
func (r *mockRequest) ContentLength() int64                { return 0 }
func (r *mockRequest) Cookie(string) (*http.Cookie, error) { return nil, nil }
func (r *mockRequest) Cookies() []*http.Cookie             { return nil }
func (r *mockRequest) Context() context.Context            { return nil }
func (r *mockRequest) Header(string) string                { return "" }
func (r *mockRequest) HeaderValues(string) []string        { return nil }
func (r *mockRequest) Headers() http.Header                { return nil }
func (r *mockRequest) Method() string                      { return r.method }
func (r *mockRequest) Protocol() string                    { return "" }
func (r *mockRequest) Path() string                        { return r.path }
func (r *mockRequest) PathValue(k string) string           { return r.vals[k] }
func (r *mockRequest) SetPathValue(k, v string)            { r.vals[k] = v }
func (r *mockRequest) Resource() string                    { return "" }
func (r *mockRequest) SetResource(string)                  {}
func (r *mockRequest) Query(string) string                 { return "" }
func (r *mockRequest) QueryValues(string) []string         { return nil }
func (r *mockRequest) Queries() url.Values                 { return nil }
func (r *mockRequest) Request() any                        { return nil }

// mockResponse implements core.Response with minimal functionality for tests.
type mockResponse struct {
	status int
}

func (r *mockResponse) AddHeader(string, string)  {}
func (r *mockResponse) SetHeader(string, string)  {}
func (r *mockResponse) DelHeader(string)          {}
func (r *mockResponse) GetHeader(string) string   { return "" }
func (r *mockResponse) Headers() http.Header      { return nil }
func (r *mockResponse) Write([]byte) (int, error) { return 0, nil }
func (r *mockResponse) Status(status int)         { r.status = status }
func (r *mockResponse) StatusCode() int           { return r.status }
func (r *mockResponse) Response() any             { return nil }

func newMockContext(method, path string) *engine.Context {
	req := &mockRequest{
		method: method,
		path:   path,
		vals:   make(map[string]string),
	}
	resp := &mockResponse{}
	return engine.NewContext(nil, req, resp)
}

func TestRoute_MatchHandlers(t *testing.T) {
	r := router.New()
	called1, called2, called3 := false, false, false
	r.GET("/api1", func(ctx core.Context) error {
		called1 = true
		return nil
	}, func(ctx core.Context) error {
		called2 = true
		return nil
	})
	r.GET("/api2", func(ctx core.Context) error {
		called3 = true
		return nil
	})

	ctx := newMockContext(http.MethodGet, "/api1")
	ctx.Use(r.Route())

	if err := r.Route()(ctx); err != nil {
		t.Fatalf("Route returned error: %v", err)
	}

	if err := ctx.Next(); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	if !called1 || !called2 || called3 {
		t.Fatalf("expected handlers h1 and h2 to be called, but got called1=%v, called2=%v, called3=%v", called1, called2, called3)
	}
}

func TestHandle_PanicOnConflict(t *testing.T) {
	r := router.New()
	isSecond := false
	defer func() {
		if rec := recover(); rec == nil || !isSecond {
			t.Fatalf("expected panic on duplicate route registration")
		}
	}()

	// first registration should succeed
	r.GET("/conflict", func(ctx core.Context) error { return nil })

	isSecond = true
	// second registration should panic
	r.GET("/conflict", func(ctx core.Context) error { return nil })
}

func TestNotFound_DefaultHandler(t *testing.T) {
	r := router.New()
	ctx := newMockContext(http.MethodGet, "/nope")

	if err := r.Route()(ctx); err != nil {
		t.Fatalf("Route returned error: %v", err)
	}
	if err := ctx.Next(); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	statusCode := ctx.Response().StatusCode()
	if statusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, statusCode)
	}
}

func TestNotFound_CustomHandler(t *testing.T) {
	r := router.New()
	r.NotFoundHandler = func(ctx core.Context) error {
		ctx.Status(http.StatusTeapot)
		ctx.Abort()
		return nil
	}
	ctx := newMockContext(http.MethodGet, "/nope")

	if err := r.Route()(ctx); err != nil {
		t.Fatalf("Route returned error: %v", err)
	}
	if err := ctx.Next(); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	statusCode := ctx.Response().StatusCode()
	if statusCode != http.StatusTeapot {
		t.Fatalf("expected status %d, got %d", http.StatusTeapot, statusCode)
	}
}
