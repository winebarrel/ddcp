PREFIX=/usr/local
VERSION=`git tag | tail -n 1`
GOOS=`go env GOOS`
GOARCH=`go env GOARCH`

ifdef GOPATH
  RUNTIME_GOPATH=$(GOPATH):`pwd`
else
  RUNTIME_GOPATH=`pwd`
endif

ddcp:	main.go src/ddcp/optparse.go src/ddcp/ddcp.go
	GOPATH=$(RUNTIME_GOPATH) go build -o ddcp main.go

install: ddcp
	install -m 755 ddcp $(DESTDIR)$(PREFIX)/bin/

clean:
	rm -f ddcp *.gz

package: clean ddcp
	gzip -c ddcp > ddcp-$(VERSION)-$(GOOS)-$(GOARCH).gz

deb:
	dpkg-buildpackage -us -uc
