GOC=go build
FETCHLIBS=go get

BUILDDIR=$(CURDIR)/build

SRCDIR=src/
OUTPUT=fsmq

GOBINDIR=$(BUILDDIR)/bin
GOPATHDIR=$(BUILDDIR)/golibs

INSTALL=install
INSTALL_BIN=$(INSTALL) -m755
INSTALL_LIB=$(INSTALL) -m644
INSTALL_CONF=$(INSTALL) -m400

PREFIX?=$(DESTDIR)/usr
BINDIR?=$(PREFIX)/bin

all: fsmq

fsmq: Makefile src/main.go
	mkdir -p $(GOPATHDIR) && \
	mkdir -p $(GOBINDIR) && \
	export GOPATH=$(GOPATHDIR) && \
	export GOBIN=$(GOBINDIR) && \
	cd $(SRCDIR) && \
	$(FETCHLIBS) && \
	$(GOC) -o $(OUTPUT)

install:
	mkdir -p $(BINDIR)
	$(INSTALL_BIN) src/$(OUTPUT) $(BINDIR)/

clean:
	rm -fr $(BUILDDIR)
	rm -f src/$(OUTPUT)
