all: 
	cd ../src
	CGO_ENABLED=0 go build . -o ./vmman4.exe

prepare:
	sudo /opt/bin/install_golang.sh 1.22.5 amd64
	sudo apk update --no-cache

install:
	test -d $(DESTDIR)$(PREFIX)/bin || mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -pm 755 $(BUILDDIR)/vmman4.exe $(DESTDIR)$(PREFIX)/bin/vmman4
	strip $(DESTDIR)$(PREFIX)/bin/vmman4
