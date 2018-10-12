<p align="center">
  <img src="https://suyashkumar.com/assets/img/lock.png" width="64">
  <h3 align="center">ssl-proxy</h3>
  <p align="center">A simple SSL reverse proxy to proxy SSL traffic to a non-SSL server with automatically generated certificates<p>
  <p align="center"><a href="https://goreportcard.com/report/github.com/suyashkumar/ssl-proxy"><img src="https://goreportcard.com/report/github.com/suyashkumar/ssl-proxy" alt=""></a> 
    <a href="https://godoc.org/github.com/suyashkumar/ssl-proxy"><img src="https://godoc.org/github.com/suyashkumar/ssl-proxy?status.svg" alt=""></a> 
  </p>
</p>


A simple single-binary SSL reverse proxy that automatically serves traffic over HTTPS and proxies it to any non-HTTPS server running on another port. `ssl-proxy` will auto-generate valid SSL certificates for your domain (from LetsEncrypt) if none are provided to it, and can also auto-generate self-signed certificates on the fly if needed for private use cases (useful for things like `jupyter` notebooks on a VM). Usage is always a simple one-liner of the form:
```sh
ssl-proxy -from 0.0.0.0:4430 -to 127.0.0.1:8000
```
or
```sh
ssl-proxy -from 0.0.0.0:4430 -to 127.0.0.1:8000 -domain=mydomain.com
```
If you want to generate and serve real SSL certificates for `mydomain.com`

## Usage
### Auto-generate and serve self-signed certificates
```sh
ssl-proxy -from 0.0.0.0:4430 -to 127.0.0.1:8000
```
This will immediately generate self-signed certificates and being proxying HTTPS traffic from https://0.0.0.0:4430 to http://127.0.0.1:8000. No need to ever call openssl. It will print the SHA256 fingerprint of the cert being used for you to perform manual certificate verification in the browser if you would like (before you "trust" the cert).

I know `nginx` is often used for stuff like this, but I got tired of dealing with the boilerplate and wanted to explore something fun. So I ended up throwing this together. 

### Provide your own certs
```sh
ssl-proxy -cert cert.pem -key myKey.pem -from 0.0.0.0:4430 -to 127.0.0.1:8000
```
You can provide your own existing certs, of course. Jenkins still has issues serving the fullchain certs from letsencrypt properly, so this tool has come in handy for me there. 

## Installation
Simply download and uncompress the proper prebuilt binary for your system from the [releases tab](https://github.com/suyashkumar/ssl-proxy/releases/). Then, add the binary to your path or start using it locally (`./ssl-proxy`).

If you're using `wget`, you can fetch and uncompress the proper binary in one command:
```sh
wget -qO- $BINARY_RELEASE_LINK | tar xvz
```

### Build from source
You must have Golang installed on your system along with `make`. Then simply clone the repository and run `make`.
