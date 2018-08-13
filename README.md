# ssl-proxy
A simple Golang SSL reverse proxy that serves traffic over HTTPS and proxies it to any other web server you might be running. `ssl-proxy` will auto-generate self-signed certificates for you if none are provided to it (useful for things like `jupyter` notebooks in a pinch). Usage is simple:

```sh
ssl-proxy -from 0.0.0.0:4430 -to 127.0.0.1:8000
```
This will immediately generate self-signed certificates and being proxying HTTPS traffic from https://0.0.0.0:4430 to http://127.0.0.1:8000. No need to ever call openssl. It will print the SHA256 fingerprint of the cert being used for you to perform manual certificate verification in the browser if you would like (before you "trust" the cert).

I know `nginx` is often used for stuff like this, but I got tired of dealing with the boilerplate and wanted to explore something fun. So I ended up throwing this together. 

## Provide your own certs
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
