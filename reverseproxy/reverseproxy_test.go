package reverseproxy

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBuild_AddHeaders tests that Build's returned ReverseProxy Director adds the proper request headers
func TestBuild_AddHeaders(t *testing.T) {
	u, err := url.Parse("http://172.16.200.81:80")
	assert.Nil(t, err, "error should be nil")
	proxy := Build(u)
	assert.NotNil(t, proxy, "Build should not return nil")

	req := httptest.NewRequest("GET", "/test", nil)
	proxy.Director(req)

	// Check that headers were added to req
	assert.Equal(t, req.Header.Get(http.CanonicalHeaderKey("X-Forwarded-Proto")), "https",
		"X-Forwarded-Proto should be present")
	assert.Equal(t, req.Header.Get(http.CanonicalHeaderKey("X-Forwarded-Port")), "443",
		"X-Forwarded-Port should be present")

}

func TestNewDirector(t *testing.T) {
	u, err := url.Parse("http://172.16.200.81:80")
	assert.Nil(t, err, "error should be nil")
	director := newDirector(u, nil)

	defaultProxy := httputil.NewSingleHostReverseProxy(u)
	defaultDirector := defaultProxy.Director

	expectedReq := httptest.NewRequest("GET", "/test", nil)
	testReq := httptest.NewRequest("GET", "/test", nil)

	defaultDirector(expectedReq)
	director(testReq)

	assert.EqualValues(t, expectedReq, testReq,
		"default proxy and package directors should modify the request in the same way")
	// TODO: add more test cases
}
