# Go HTTP Proxy

Go HTTP Proxy is a simple proxy server written in Go. It provides the ability to proxy HTTP and HTTPS traffic leveraging the HTTP CONNECT Method.

This tool is useful when you need to debug or inspect network requests, or to test your client code in presence of proxies.

## Features

- Supports both HTTP and HTTPS. Forwards HTTPS traffic using the HTTP CONNECT method.
- Provides an option to increase verbosity to log detailed information about the traffic.
- Supports custom TLS certificates (`-cert` and `-key` flags).

## Limitations

- This tool is primarily for debugging and development purposes. It might not be suited for production-level traffic or security demands.
- The proxy does not modify or interfere with the content of the traffic, it merely passes it along. This means it will not remove certain headers or modify request/response bodies.

## Usage

By default, the proxy will start an HTTP server on port 8888.

```bash
./go-http-proxy --help
  -cert string
        Path to TLS certificate. If not provided and https is set, a self-signed certificate will be generated and saved to cert.pem in the current directory.
  -https
        Use HTTPS for proxy
  -key string
        Path to TLS key. If not provided and https is set, a self-signed certificate will be generated.
  -port string
        Port to listen on (default "8888")
  -v    Increase verbosity
```

## Get Started

First, download the latest binary for your platform from the GitHub releases page. For example for Linux x64 or AMD64

```bash
curl -LOs https://github.com/kitos9112/go-http-proxy/releases/latest/download/go-http-proxy_Linux_x86_64.tar.gz
tar -xzf go-http-proxy_Linux_x86_64.tar.gz
./go-http-proxy -v &
echo "All set! Now you can use the proxy on  http://localhost:8888"
curl -k --proxy http://localhost:8888 https://google.com -I
```

### Setting up a HTTPs Proxy

To start a HTTPS server, use the `-https` flag and optionally provide paths to your certificate and key files.

If `key` and `cert` paths are not passed, the program will generate a self-signed certificate and key with a CNAME and DNS SAN of `localhost` for easy DNS resolution during local tests. The x509 cert will be written to disk in the current directory as `cert.pem` - No key will be exported though.

You'd either have to add this cert to the `ca-certificates` store on your local machine, or avoid verifying the cert at all from the client.

An example on how to use this proxy with a self-signed certificate and then importing to the local CA store on Ubuntu:

```bash
./go-http-proxy -https -port 4443 -v &
sleep 2
sudo cp cert.pem /usr/local/share/ca-certificates/self-signed.crt
sudo update-ca-certificates && sudo update-ca-certificates --fresh
curl -k --proxy https://localhost:4443 https://google.com -I
```

Most client libraries have a standard way (Common Gateway Interface specification) to pass HTTP proxy information via environment variables. For example, for HTTP traffic you should set the `HTTP_PROXY` or `http_proxy` env. variable whereas for HTTPS traffic you should set the `HTTPS_PROXY` or `https_proxy` env. variable.
