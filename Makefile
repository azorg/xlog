# File: "Makefile"

PRJ="github.com/azorg/xlog"

.PHONY: all distclean fmt test

all: fmt test

distclean:
	@rm go.mod

go.mod:
	@go mod init $(PRJ)

fmt: go.mod
	@go fmt

test: go.mod
	@go test

# EOF: "Makefile"
