GO		= go
# -s removes symbol table and -ldflags -w debugging symbols
LDFLAGS		= -trimpath -ldflags "-s -w"
GOOS		= linux
GOARCH		= amd64

.PHONY: all clean analysis obsd test

# Defaults Linux
all: build

build:
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=$(GOOS) \
		$(GO) build $(LDFLAGS)

lint:
	gosec ./...
	staticcheck ./...
	go vet ./...

debug:
	CGO_ENABLED=1 $(GO) build $(LDFLAGS)

obsd:
	GOOS=openbsd $(GO) build $(LDFLAGS) -o streamjury_obsd

test:
	go test ./...

clean:
	rm -f streamjury streamjury_obsd

