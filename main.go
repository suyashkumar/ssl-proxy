package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"crypto/tls"

	"github.com/eastlondoner/tailscale-ssl-proxy/reverseproxy"

	"tailscale.com/client/tailscale"
)

var (
	to           = flag.String("to", "http://localhost:8080", "the address and port for which to proxy requests to")
	fromURL      = flag.String("from", ":443", "the tcp address and port this proxy should listen for requests on")
	redirectHTTP = flag.String("redirectHTTP", ":80", "the tcp address and port this proxy should listen for http->https request redirects. Set to 'off' to disable http->https redirect")
)

const (
	HTTPSPrefix     = "https://"
	HTTPPrefix      = "http://"
)

func main() {
	flag.Parse()

	// Ensure the from URL is in the right form
	if !strings.Contains(*fromURL, ":") {
		log.Fatal("Unable to parse 'from' url (did you specify a port?)")
	}
	port := (*fromURL)[strings.LastIndex(*fromURL, ":")+1:]


	// Ensure the to URL is in the right form
	if strings.HasPrefix(*to, ":") {
		*to = (*to)[1:]
		log.Println("Assuming -to URL a port number")
	}

	toPort, err := strconv.Atoi(*to)
    if err == nil && toPort > 0 {
		// to is just a port number
		*to = "localhost:" + *to
		log.Println("Assuming -to should use localhost")
    }

	if !strings.HasPrefix(*to, HTTPPrefix) && !strings.HasPrefix(*to, HTTPSPrefix) {
		*to = HTTPPrefix + *to
		log.Println("Assuming -to URL is using http://")
	}

	// Parse toURL as a URL
	toURL, err := url.Parse(*to)
	if err != nil {
		log.Fatal("Unable to parse 'to' url: ", err)
	}

	// Create main TCP listener
	ln, err := net.Listen("tcp", *fromURL)
	if err != nil {
		log.Printf("Unable to listen on %s", *fromURL)
		log.Fatal(err)
	} else {
		log.Printf("Listening on %s", ln.Addr().String())
	}

	// Redirect http requests on port 80 to TLS port using https
	if *redirectHTTP != "off" {
		// Redirect to fromURL by default, unless a domain is specified--in that case, redirect using the public facing
		// domain

		// Create redirect TCP listener
		redirectListener, err := net.Listen("tcp", *redirectHTTP)
		if err != nil {
			log.Printf("Unable to listen on %s", *redirectHTTP)
			log.Fatal(err)
		} else {
			log.Printf("Listening on %s", ln.Addr().String())
		}

		redirectTLS := func(w http.ResponseWriter, r *http.Request) {
			redirectURL := "https://"+r.URL.Hostname()
			if port != "443" {
				redirectURL = redirectURL+":"+port
			}
			redirectURL = redirectURL+r.RequestURI
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
		}

		go func() {
			log.Printf("Redirecting http requests on %s to https on port :%s", *redirectHTTP, port)
			err := http.Serve(redirectListener, http.HandlerFunc(redirectTLS))
			if err != nil {
				log.Println("HTTP redirection server failure")
				log.Fatal(err)
			}
		}()
	}

	// Setup reverse proxy ServeMux
	p := reverseproxy.Build(toURL)
	mux := http.NewServeMux()
	mux.Handle("/", p)

	log.Printf(green("Proxying calls from https://%s (SSL/TLS) to %s"), *fromURL, toURL)
	
	if port != "443" {
		log.Println("WARN: You must serve on port :443 for the LetsEncrypt certs that tailscale uses to be valid")
	}

	s := &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: tailscale.GetCertificate,
		},
	}
	s.Handler = mux
	log.Fatal(s.ServeTLS(ln,"",""))
}

// green takes an input string and returns it with the proper ANSI escape codes to render it green-colored
// in a supported terminal.
// TODO: if more colors used in the future, generalize or pull in an external pkg
func green(in string) string {
	return fmt.Sprintf("\033[0;32m%s\033[0;0m", in)
}
