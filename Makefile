.PHONY: install-deps test build

install-deps:
	go get github.com/spf13/cobra
	go get github.com/onsi/ginkgo
	go get github.com/onsi/gomega

test:
	./test.sh

build:
	go build
