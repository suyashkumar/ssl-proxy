<p align="center">
  <h3 align="center">tailscale-ssl-proxy</h3>
  <p align="center">Simple single-command SSL reverse proxy for Tailscale<p>
</p>

A handy and simple way to add Tailscale SSL support to your locally running thing --be it your personal jupyter notebook, nodejs app or whatever. `tailscale-ssl-proxy` uses the official (Tailscale go package)[https://pkg.go.dev/tailscale.com] to get trusted LetsEncrypt SSL certs and then proxies HTTPS traffic to your existing HTTP server in a single command. `tailscale-ssl-proxy` also redirects unencrypted HTTP traffic on port 80 to HTTPS.

## Usage

### Simple
```sh
tailscale-ssl-proxy
```
This will immediately fetch, real LetsEncrypt certificates for the machine's Tailscale address.

### Specify the 
```sh
tailscale-ssl-proxy -from 0.0.0.0:4430 -to 127.0.0.1:8000
```
This will immediately generate self-signed certificates and begin proxying HTTPS traffic from https://0.0.0.0:4430 to http://127.0.0.1:8080. No need to ever call openssl. It will print the SHA256 fingerprint of the cert being used for you to perform manual certificate verification in the browser if you would like (before you "trust" the cert).

I know `nginx` is often used for stuff like this, but I got tired of dealing with the boilerplate and wanted to explore something fun. So I ended up throwing this together. 

### Redirect HTTP -> HTTPS
Simply include the `-redirectHTTP` flag when running the program.

### Build from source 
#### Build from source using Docker
You can build `tailscale-ssl-proxy` for all platforms quickly using the included Docker configurations.

If you have `docker-compose` installed:
```sh
docker build . -t tailscale-ssl-proxy_build-release
docker-compose -f docker-compose.build.yml up
```
will build linux, osx, and darwin binaries (x86) and place them in a `build/` folder in your current working directory.
#### Build from source locally
You must have Golang installed on your system along with `make`. Then simply clone the repository and run `make`. 

## Attribution
Forked from <a href="https://github.com/suyashkumar/ssl-proxy">ssl-proxy by Suyash Kumar</a>
Icons made by <a href="https://www.flaticon.com/authors/those-icons" title="Those Icons">Those Icons</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a>
