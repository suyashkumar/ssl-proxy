package reverseproxy

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
)

// TestBuild_AddHeaders tests that Build's returned ReverseProxy Director adds the proper request headers
func TestBuild_AddHeaders(t *testing.T) {
	u, err := url.Parse("http://127.0.0.1")
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	proxy := Build(u)
	if proxy == nil {
		t.Fatal("got nil, want non-nil proxy")
	}

	req := httptest.NewRequest("GET", "/test", nil)
	proxy.Director(req)

	// Check that headers were added to req
	if got := req.Header.Get(http.CanonicalHeaderKey("X-Forwarded-Proto")); got != "https" {
		t.Errorf("X-Forwarded-Proto: got %q, want %q", got, "https")
	}
	if got := req.Header.Get(http.CanonicalHeaderKey("X-Forwarded-Port")); got != "443" {
		t.Errorf("X-Forwarded-Port: got %q, want %q", got, "443")
	}
}

func TestNewDirector(t *testing.T) {
	u, err := url.Parse("http://127.0.0.1")
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	director := newDirector(u, nil)

	defaultProxy := httputil.NewSingleHostReverseProxy(u)
	defaultDirector := defaultProxy.Director

	expectedReq := httptest.NewRequest("GET", "/test", nil)
	testReq := httptest.NewRequest("GET", "/test", nil)

	defaultDirector(expectedReq)
	director(testReq)

	// Compare relevant fields of the requests
	if got, want := testReq.URL.String(), expectedReq.URL.String(); got != want {
		t.Errorf("URL: got %q, want %q", got, want)
	}
	if got, want := testReq.Host, expectedReq.Host; got != want {
		t.Errorf("Host: got %q, want %q", got, want)
	}
	if got, want := testReq.RemoteAddr, expectedReq.RemoteAddr; got != want {
		t.Errorf("RemoteAddr: got %q, want %q", got, want)
	}
	// TODO: add more test cases
}
