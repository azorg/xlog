# File: "Makefile"

PRJ="github.com/azorg/xlog"

.PHONY: all distclean fmt test

all: fmt test

distclean:
	@rm -f go.mod
	@rm -f go.sum
	@#sudo rm -rf go/pkg
	@rm -rf vendor
	@go clean -modcache
	
go.mod:
	@go mod init $(PRJ)

go.sum: go.mod Makefile
	@#go get golang.org/x/exp/slog@master # experimental slog (go <1.21)
	@touch go.sum

tidy: go.mod
	@go mod tidy # automatic update go.sum

vendor: go.sum
	@go mod vendor

fmt: go.mod go.sum
	@go fmt

test: go.mod go.sum
	@go test

# EOF: "Makefile"
