.PHONY: build

PROG = airgeek-rpi

build:
	go build -o $(PROG)

install: build
	install -DT -m 755 $(PROG) $(DESTDIR)$(PREFIX)/bin/$(PROG)
	install -DT -m 644 $(PROG).service /etc/systemd/system/$(PROG).service
	install -DT -m 755 $(PROG)-launch.sh $(DESTDIR)$(PREFIX)/bin/$(PROG)-launch.sh
