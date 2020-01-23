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
	@go build -o ${GOPATH}/bin/release-version -i ${BUILD_FLAGS} main.go