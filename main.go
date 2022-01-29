package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"crypto/tls"

	"github.com/suyashkumar/ssl-proxy/gen"
	"github.com/suyashkumar/ssl-proxy/reverseproxy"
	"golang.org/x/crypto/acme/autocert"

	"tailscale.com/client/tailscale"
)

var (
	to           = flag.String("to", "http://127.0.0.1:80", "the address and port for which to proxy requests to")
	fromURL      = flag.String("from", "127.0.0.1:4430", "the tcp address and port this proxy should listen for requests on")
	redirectHTTP = flag.Bool("redirectHTTP", true, "redirect http requests from port 80 to https (enabled by default)")
)

const (
	HTTPSPrefix     = "https://"
	HTTPPrefix      = "http://"
)

func main() {
	flag.Parse()

	// Ensure the to URL is in the right form
	if !strings.HasPrefix(*to, HTTPPrefix) && !strings.HasPrefix(*to, HTTPSPrefix) {
		*to = HTTPPrefix + *to
		log.Println("Assuming -to URL is using http://")
	}

	// Parse toURL as a URL
	toURL, err := url.Parse(*to)
	if err != nil {
		log.Fatal("Unable to parse 'to' url: ", err)
	}

	// Setup reverse proxy ServeMux
	p := reverseproxy.Build(toURL)
	mux := http.NewServeMux()
	mux.Handle("/", p)

	log.Printf(green("Proxying calls from https://%s (SSL/TLS) to %s"), *fromURL, toURL)

	// Redirect http requests on port 80 to TLS port using https
	if *redirectHTTP {
		// Redirect to fromURL by default, unless a domain is specified--in that case, redirect using the public facing
		// domain
		redirectURL := *fromURL
		if validDomain {
			redirectURL = *domain
		}
		redirectTLS := func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.URL.Host+r.RequestURI, http.StatusMovedPermanently)
		}
		go func() {
			log.Println(
				fmt.Sprintf("Also redirecting https requests on port 80 to https requests on %s", redirectURL))
			err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS))
			if err != nil {
				log.Println("HTTP redirection server failure")
				log.Println(err)
			}
		}()
	}

	
	if !strings.HasSuffix(*fromURL, ":443") {
		log.Println("WARN: Right now, you must serve on port :443 for the LetsEncrypt certs that tailscale uses to be valid")
	}

	s := &http.Server{
		Addr:      *fromURL,
		TLSConfig: &tls.Config{
			GetCertificate: tailscale.GetCertificate,
		},
	}
	s.Handler = mux
	log.Fatal(s.ListenAndServeTLS("", ""))
}

// green takes an input string and returns it with the proper ANSI escape codes to render it green-colored
// in a supported terminal.
// TODO: if more colors used in the future, generalize or pull in an external pkg
func green(in string) string {
	return fmt.Sprintf("\033[0;32m%s\033[0;0m", in)
}
