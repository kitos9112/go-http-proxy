package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	verbose     = flag.Bool("v", false, "Increase verbosity")
	port        = flag.String("port", "8888", "Port to listen on")
	https       = flag.Bool("https", false, "Use HTTPS for proxy")
	certificate = flag.String("cert", "", "Path to TLS certificate. If not provided and https is set, a self-signed certificate will be generated and saved to cert.pem in the current directory.")
	key         = flag.String("key", "", "Path to TLS key. If not provided and https is set, a self-signed certificate will be generated.")
)

func verbosePrintln(logger *log.Logger, msg string) {
	if *verbose {
		logger.Println(msg)
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	verbosePrintln(log.New(io.Discard, "", log.LstdFlags), "Handling CONNECT request for "+r.Host)

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	_, err := io.Copy(destination, source)
	if err != nil {
		log.Println("Error while transferring data:", err)
		return // Return without crashing the server.
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	verbosePrintln(log.New(io.Discard, "", log.LstdFlags), "Handling HTTP request for "+r.URL.String())
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil && *verbose {
		log.Println("Error while copying response body:", err)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func genTlsCert() *tls.Config {
	// Generate a new self-signed certificate.
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber: big.NewInt(2024),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		DNSNames:              []string{"localhost"},
		NotBefore:             time.Now().AddDate(0, 0, -1),         // Valid from one day before
		NotAfter:              time.Now().Add(time.Hour * 24 * 365), // Valid for one year.
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// Write the certificate to file.
	if err := os.WriteFile("cert.pem", certPEM, 0644); err != nil {
		log.Fatalf("Failed to write certificate: %v", err)
	}

	cert, _ := tls.X509KeyPair(certPEM, keyPEM)
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func main() {
	flag.Parse()

	server := &http.Server{
		Addr: ":" + *port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
	}

	if *https {
		var tlsConfig *tls.Config

		if *certificate != "" && *key != "" {
			cert, err := tls.LoadX509KeyPair(*certificate, *key)
			if err != nil {
				log.Fatalf("Failed to load certificate and key: %v", err)
			}

			tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		} else {
			tlsConfig = genTlsCert()
		}

		server.TLSConfig = tlsConfig

		log.Println("Starting HTTPS server on", server.Addr)
		log.Fatal(server.ListenAndServeTLS("", ""))

	} else {
		log.Println("Starting HTTP server on", server.Addr)
		log.Fatal(server.ListenAndServe())
	}
}
