# CGO_ENABLED=0 == static by default
GO		= go
# -s removes symbol table and -ldflags -w debugging symbols
LDFLAGS		= -trimpath -ldflags "-s -w"
GOOS		= linux
GOARCH		= amd64
# XX For future release -gcflags=all=-d=checkptr -run=Rights syscall

.PHONY: all analysis obsd test

# Defaults Linux
all:
	CGO_ENABLED=0 $(GO) build $(LDFLAGS)
debug:
	CGO_ENABLED=1 $(GO) build $(LDFLAGS)
obsd:
	GOOS=openbsd $(GO) build $(LDFLAGS) -o streamjury_obsd
test:
	go test ./...
