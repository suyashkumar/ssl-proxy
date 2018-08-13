# ssl-proxy
A simple Golang SSL reverse proxy that serves traffic over HTTPS and proxies it to any other web server you might be running. Overall, a simple way to add SSL/TLS to a web service. `ssl-proxy` will auto-generate self-signed certificates for you if none are provided to it (useful for things like `jupyter` notebooks in a pinch). 

I know `nginx` is often used for stuff like this, but I got tired of dealing with the boilerplate and wanted to explore something fun. So I ended up throwing this together. 

## Quick Start
```sh
ssl-proxy -from 0.0.0.0:4430 -to http://127.0.0.1:8000
```
This will immediately generate self-signed certificates and being proxying HTTPS traffic from `0.0.0.0:4430` to `http://127.0.0.1:8000`. No need to ever call `openssl`. 

It will output the SHA256 fingerprint for you to perform manual certificate verification in the browser if needed (before you "trust" the cert). 

## Provide your own certs
```sh
ssl-proxy -cert cert.pem -key myKey.pem -from 0.0.0.0:4430 -to http://127.0.0.1:8000
```
You can provide your own existing certs, of course. Jenkins still has issues serving the fullchain certs from letsencrypt properly, so this tool has come in handy for me there. 
