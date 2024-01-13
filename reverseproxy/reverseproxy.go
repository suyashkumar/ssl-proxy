package reverseproxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Build initializes and returns a new ReverseProxy instance suitable for SSL proxying
func Build(toURL *url.URL, ipFilter string) *httputil.ReverseProxy {
	localProxy := &httputil.ReverseProxy{}
	addProxyHeaders := func(req *http.Request) {
		req.Header.Set(http.CanonicalHeaderKey("X-Forwarded-Proto"), "https")
		req.Header.Set(http.CanonicalHeaderKey("X-Forwarded-Port"), "443") // TODO: inherit another port if needed
	}
	localProxy.Director = newDirector(toURL, ipFilter, addProxyHeaders)

	return localProxy
}

// newDirector creates a base director that should be exactly what http.NewSingleHostReverseProxy() creates, but allows
// for the caller to supply and extraDirector function to decorate to request to the downstream server
func newDirector(target *url.URL, ipFilter string, extraDirector func(*http.Request)) func(*http.Request) {
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		remoteIp, _, _ := net.SplitHostPort(req.RemoteAddr)
		if ipFilter == "" || remoteIp == ipFilter {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}

			if extraDirector != nil {
				extraDirector(req)
			}
		} else {
			// send to black hole
			req.URL.Host = "127.0.0.1:0"
			req.URL.Scheme = "http"
		}
	}
}

// singleJoiningSlash is a utility function that adds a single slash to a URL where appropriate, copied from
// the httputil package
// TODO: add test to ensure behavior does not diverge from httputil's implementation, as per Rob Pike's proverbs
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
