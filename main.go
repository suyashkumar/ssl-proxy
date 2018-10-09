package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"strings"

	"github.com/suyashkumar/ssl-proxy/gen"
	"golang.org/x/crypto/acme/autocert"
)

var (
	to       = flag.String("to", "http://127.0.0.1:80", "the address and port for which to proxy requests to")
	fromURL  = flag.String("from", "127.0.0.1:4430", "the tcp address and port this proxy should listen for requests on")
	certFile = flag.String("cert", "", "path to a tls certificate file. If not provided, ssl-proxy will generate one for you in ~/.ssl-proxy/")
	keyFile  = flag.String("key", "", "path to a private key file. If not provided, ssl-proxy will generate one for you in ~/.ssl-proxy/")
	domain = flag.String("domain", "", "domain to mint letsencrypt certificates for. Usage of this parameter implies acceptance of the LetsEncrypt terms of service.")
)

const (
	DefaultCertFile = "cert.pem"
	DefaultKeyFile  = "key.pem"
	HTTPSPrefix     = "https://"
	HTTPPrefix      = "http://"
)

func main() {
	flag.Parse()

	if *certFile == "" || *keyFile == "" {
		// Use default file paths
		*certFile = DefaultCertFile
		*keyFile = DefaultKeyFile

		log.Printf("No existing cert or key specified, generating some self-signed certs for use (%s, %s)\n", *certFile, *keyFile)

		// Generate new keys
		certBuf, keyBuf, fingerprint, err := gen.Keys(365 * 24 * time.Hour)
		if err != nil {
			log.Fatal("Error generating default keys", err)
		}

		certOut, err := os.Create(*certFile)
		if err != nil {
			log.Fatal("Unable to create cert file", err)
		}
		certOut.Write(certBuf.Bytes())

		keyOut, err := os.OpenFile(*keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatal("Unable to create the key file", err)
		}
		keyOut.Write(keyBuf.Bytes())

		log.Printf("SHA256 Fingerprint: % X", fingerprint)

	}

	// Ensure the to URL is in the right form
	if !strings.HasPrefix(*to, HTTPPrefix) && !strings.HasPrefix(*to, HTTPSPrefix) {
		*to = HTTPPrefix + *to
		log.Println("Assuming -to URL is using http://")
	}

	toURL, err := url.Parse(*to)
	if err != nil {
		log.Fatal("Unable to parse 'to' url: ", err)
	}

	// Setup ServeMux
	localProxy := httputil.NewSingleHostReverseProxy(toURL)
	mux := http.NewServeMux()
	mux.Handle("/", localProxy)

	if *domain != "" {
		// Domain is present, use autocert
		// TODO: validate domain (though, autocert may do this)
		log.Fatal(http.Serve(autocert.NewListener(*domain), mux))
	} else {
		log.Printf("Proxying calls from https://%s (SSL/TLS) to %s", *fromURL, toURL)
		log.Fatal(http.ListenAndServeTLS(*fromURL, *certFile, *keyFile, mux))
	}
}
