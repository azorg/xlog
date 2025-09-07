PRJ = "github.com/azorg/xlog"

# Version, git hash
MAJOR := 2
MINOR := 0
BUILD := 0
VERSION := $(MAJOR).$(MINOR).$(BUILD)
GIT_HASH := `git rev-parse HEAD | head -c 7`
BUILD_TIME := `date '+%Y.%m.%d_%H:%M'`

GIT_MESSAGE = "fix: смотри ChangeLog"

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.GitHash=$(GIT_HASH) -X=main.BuildTime=$(BUILD_TIME)"
GCFLAGS= 

# Private repository (xlog)
#GIT_REPO := github.com/azorg
#XLOG := $(GIT_REPO)/xlog
#export GOPRIVATE=$(XLOG),$(XLOG)/signal

.PHONY: all fmt simplify vet test clean distclean tidy vendor doc commit

all: go.mod tidy fmt simplify vet doc xlogscan #test

fmt:
	go fmt

simplify:
	gofmt -l -w -s *.go

vet:
	go vet

test:
	go test

clean:
	rm -rf ./logs
	rm -f ./xlogscan
	@#rm -f README.*

#prepare-git:
#	@echo ">>> prepare private ssh access"
#	@git config --global url."git@$(GIT_REPO):".insteadOf "https://$(GIT_REPO)/"

distclean: clean
	rm -f go.mod
	rm -f go.sum
	@#sudo rm -rf go/pkg
	rm -rf vendor
	go clean -modcache

go.mod:
	@echo ">>> create go.mod"
	@go mod init $(PRJ)
	@touch go.mod

go.sum: go.mod Makefile tidy
	@echo ">>> create go.sum"
	@#go get golang.org/x/exp/slog@v0.0.0-20240904232852-e7e105dedf7e # experimental slog (go <=1.20)
	@go get gopkg.in/natefinch/lumberjack.v2 # Lumberjack as log rotate
	@go get github.com/gofrs/uuid # UUID (v7)
	@go get github.com/sigurn/crc16 # CRC16
	@#go get github.com/sigurn/crc8 # CRC8
	@touch go.sum

tidy: go.mod
	@echo ">>> automatic update go.sum by tidy"
	@go mod tidy # automatic update go.sum

vendor: go.sum
	@echo ">>> create vendor"
	@go mod vendor

doc: README.txt README.md

README.txt: *.go
	go doc -all > README.txt

README.md: *.go ~/go/bin/gomarkdoc
	~/go/bin/gomarkdoc -o README.md

~/go/bin/gomarkdoc:
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

commit: doc
	git add .
	git commit -am $(GIT_MESSAGE)
	git push

xlogscan: *.go cmd/xlogscan/*.go
	@echo ">>> build $@"
	@go build $(LDFLAGS) $(GCFLAGS) -o . $(PRJ)/cmd/xlogscan/

