PREFIX ?= /usr
BINARY_NAME=shlog

all: build

deps:
	go get atomicgo.dev/keyboard/keys
	go get github.com/pterm/pterm

build:
	go build -o ${BINARY_NAME} src/main.go

run:
	go build -o ${BINARY_NAME} src/main.go
	./${BINARY_NAME}

install:
	@install -Dm755 shlog $(DESTDIR)$(PREFIX)/bin/shlog

uninstall:
	@rm -f $(DESTDIR)$(PREFIX)/bin/shlog

clean:
	go clean
	rm ${BINARY_NAME}
