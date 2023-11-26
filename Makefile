# File: "Makefile"

PRJ="github.com/azorg/xlog"

GIT_MESSAGE = "auto commit"

.PHONY: all help distclean commit tidy vendor fmt test

all: fmt test

help:
	@echo "make fmt       - format Go sources"
	@echo "make test      - run test"
	@echo "make all       - fmt + test"
	@echo "make distclean - full clean (go.mod, go.sum)"

distclean:
	@rm -f go.mod
	@rm -f go.sum
	@#sudo rm -rf go/pkg
	@rm -rf vendor
	@go clean -modcache
	
commit:
	git add .
	git commit -am $(GIT_MESSAGE)
	git push

go.mod:
	@go mod init $(PRJ)

go.sum: go.mod Makefile
	@#go get golang.org/x/exp/slog # experimental slog (go <1.21)
	@touch go.sum

tidy: go.mod
	@go mod tidy

vendor: go.sum
	@go mod vendor

fmt: go.mod go.sum
	@go fmt

test: go.mod go.sum
	@go test

# EOF: "Makefile"
