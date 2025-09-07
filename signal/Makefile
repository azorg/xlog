.PHONY: all fmt simplify vet test clean doc

all: fmt simplify vet doc test

fmt:
	go fmt

simplify:
	gofmt -l -w -s *.go

vet:
	go vet

test:
	go test

clean:
	@#rm -f README.*

doc: README.txt README.md

README.txt: *.go
	go doc -all > README.txt

README.md: *.go ~/go/bin/gomarkdoc
	~/go/bin/gomarkdoc -o README.md

~/go/bin/gomarkdoc:
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

