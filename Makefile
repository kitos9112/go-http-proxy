VERSION = $(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"
OSARCH=$(shell go env GOHOSTOS)-$(shell go env GOHOSTARCH)

SIMPLEPROXYTOOL=\
	go-hhtp-proxy-darwin-amd64 \
	go-http-proxy-linux-amd64 \
	go-http-proxy-windows-amd64.exe

proxytool: go-http-proxy-$(OSARCH)

$(SIMPLEPROXYTOOL):
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$(subst .exe,,$@))) go build $(LDFLAGS) -o $@ ./$<

%-$(VERSION).zip: %.exe
	rm -f $@
	zip $@ $<

%-$(VERSION).zip: %
	rm -f $@
	zip $@ $<

clean:
	rm -f go-http-proxy-*


release:
	$(foreach bin,$(SIMPLEPROXYTOOL),$(subst .exe,,$(bin))-$(VERSION).zip)

rm:
	$(foreach bin,$(SIMPLEPROXYTOOL),$(subst .exe,,$(bin))-$(VERSION).zip)


build: proxytool $(SIMPLEPROXYTOOL) clean rm

.PHONY: build
all: build
