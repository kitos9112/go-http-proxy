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

First, download the latest binary for your platform from the GitHub releases page. Then you can run the proxy with:

By default, the proxy will start an HTTP server on port 8888. By default, the proxy will start an HTTP server on port 8888.

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

To start an HTTPS server, use the `-https` flag and optionally provide paths to your certificate and key files.

If `key` and `cert` paths are not passed, the program will generate a self-signed certificate and key with a CNAME and DNS SAN of `localhost` for easy DNS resolution during local tests. Then, you'd either have to add this cert to the `ca-certificates` store on your machine, or avoid verifying the cert in the client.

```bash
./go-http-proxy -https
```
