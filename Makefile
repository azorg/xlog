# File: "Makefile"

PRJ="github.com/azorg/xlog"

GIT_MESSAGE = "auto commit"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# go packages
PKGS = $(PRJ)

.PHONY: all help distclean clean commit tidy vendor fmt vet doc test

all: fmt vet doc test

help:
	@echo "make fmt       - format Go sources"
	@echo "make test      - run test"
	@echo "make distclean - full clean (go.mod, go.sum)"

clean:
	rm -rf logs

distclean:
	@rm -f go.mod
	@rm -f go.sum
	@#sudo rm -rf go/pkg
	@rm -rf vendor
	@go clean -modcache
	
commit: clean fmt
	git add .
	git commit -am $(GIT_MESSAGE)
	git push

go.mod:
	@go mod init $(PRJ)
	@touch go.mod

tidy: go.mod
	@go mod tidy

go.sum: go.mod Makefile #tidy
	@go get golang.org/x/exp/slog@v0.0.0-20240904232852-e7e105dedf7e # experimental slog (go <1.21)
	@go get gopkg.in/natefinch/lumberjack.v2 # Lumberjack as log rotate
	@go get github.com/google/uuid # Google UUID
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

doc: README-RU.txt README-RU.md

README-RU.txt: *.go
	go doc -all > README-RU.txt

README-RU.md: *.go ~/go/bin/gomarkdoc
	~/go/bin/gomarkdoc -o README-RU.md

~/go/bin/gomarkdoc:
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

# EOF: "Makefile"
