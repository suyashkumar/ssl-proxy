package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/suyashkumar/ssl-proxy/gen"
)

var (
	to       = flag.String("to", "http://127.0.0.1:80", "the address and port for which to proxy requests to")
	fromURL  = flag.String("from", "127.0.0.1:4430", "the tcp address and port this proxy should listen for requests on")
	certFile = flag.String("cert", "", "path to a tls certificate file. If not provided, ssl-proxy will generate one for you in ~/.ssl-proxy/")
	keyFile  = flag.String("key", "", "path to a private key file. If not provided, ssl-proxy will generate one for you in ~/.ssl-proxy/")
)

const (
	DefaultCertFile = "cert.pem"
	DefaultKeyFile  = "key.pem"
)

func main() {
	flag.Parse()

	if *certFile == "" || *keyFile == "" {
		// Use default file paths
		*certFile = DefaultCertFile
		*keyFile = DefaultKeyFile

		log.Printf("No existing cert or key specified, generating some self-signed certs for use (%s, %s)\n", *certFile, *keyFile)

		// Generate new keys
		certBuf, keyBuf, err := gen.Keys(365 * 24 * time.Hour)
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

	}

	toURL, err := url.Parse(*to)
	if err != nil {
		log.Fatal("Unable to parse 'to' url: ", err)
	}

	localProxy := httputil.NewSingleHostReverseProxy(toURL)
	http.Handle("/", localProxy)
	log.Printf("Proxying calls from %s (SSL/TLS) to %s", *fromURL, toURL)
	log.Fatal(http.ListenAndServeTLS(*fromURL, *certFile, *keyFile, nil))
}
