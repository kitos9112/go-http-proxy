package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestCopyHeader(t *testing.T) {
	dst := make(http.Header)
	src := make(http.Header)

	src.Add("Content-Type", "text/html")
	src.Add("Content-Length", "123")

	copyHeader(dst, src)

	if dst.Get("Content-Type") != "text/html" || dst.Get("Content-Length") != "123" {
		t.Fatal("copyHeader did not correctly copy headers from src to dst")
	}
}

func TestVerbosePrintln(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.LstdFlags)

	*verbose = true

	verbosePrintln(logger, "Test message")

	if !strings.Contains(buf.String(), "Test message") {
		t.Fatal("verbosePrintln did not correctly output message to log")
	}

	*verbose = false
}
