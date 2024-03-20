#
# Tools and binaries
#
GOCMD		= go
GOTEST		=$(GOCMD) test

#
# Directories and packages
#
TEST_PKGS := $(shell go list ./...)

#
# Targets
#
.PHONY: test
test:
	$(GOTEST) $(TEST_PKGS)
.PHONY: testv
testv:
	$(GOTEST) -v $(TEST_PKGS)

clean:
	@rm -f release-version

install:
	@go build -o ${GOPATH}/bin/release-version main.go

fmt:
	goimports -local github.com/flume -w ./
