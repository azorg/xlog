# File: "Makefile"

PRJ="github.com/azorg/xlog"

GIT_MESSAGE = "auto commit"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# go packages
PKGS = $(PRJ)

.PHONY: all help distclean commit tidy vendor fmt test

all: fmt test

help:
	@echo "make fmt       - format Go sources"
	@echo "make test      - run test"
	@echo "make distclean - full clean (go.mod, go.sum)"

distclean:
	@rm -f go.mod
	@rm -f go.sum
	@#sudo rm -rf go/pkg
	@rm -rf vendor
	@go clean -modcache
	
commit: fmt
	git add .
	git commit -am $(GIT_MESSAGE)
	git push

go.mod:
	@go mod init $(PRJ)
	@touch go.mod

tidy: go.mod
	@go mod tidy

go.sum: go.mod Makefile #tidy
	@#go get golang.org/x/exp/slog # experimental slog (go <1.21)
	@go get github.com/gofrs/uuid # UUID v1-v7
	@touch go.sum

vendor: go.sum
	@go mod vendor

fmt: go.mod go.sum
	@go fmt

simplify:
	@gofmt -l -w -s $(SRC)

vet:
	@#go vet
	@go vet $(PKGS)

test: go.mod go.sum
	@go test

# EOF: "Makefile"
