.DEFAULT_GOAL = all
SHELL         = bash

skip = $(info $@: skipping, target disabled)

# Git
#
# Provide some nice to use variables for the git
# repository state
COMMIT := $(shell git rev-parse HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
SLUG   := $(shell git remote -v | grep "(fetch)" | awk '{print$$2}' | sed -E 's/^.*(\/|:)([^ ]*)\/([^ ]*)$$/\2\/\3/;s/\.git//')
OS     := $(shell uname)

# Directories
#
# All of the following directories can be
# overwritten. If this is done, it is
# only recommended to change the BUILD_DIR
# option.
BUILD_DIR     := build
RELEASE_DIR   := $(BUILD_DIR)/release
LINT_DIR      := $(BUILD_DIR)/lint
TEST_DIR      := $(BUILD_DIR)/test
IMAGE_DIR     := $(BUILD_DIR)/container
DIST_DIR      := $(BUILD_DIR)/dist
INT_DIR       := $(BUILD_DIR)/integration

$(BUILD_DIR):
	-mkdir $(BUILD_DIR)

$(RELEASE_DIR): | $(BUILD_DIR)
	-mkdir $(RELEASE_DIR)

$(LINT_DIR): | $(BUILD_DIR)
	-mkdir $(LINT_DIR)

$(TEST_DIR): | $(BUILD_DIR)
	-mkdir $(TEST_DIR)

$(IMAGE_DIR): | $(BUILD_DIR)
	-mkdir $(IMAGE_DIR)

$(DIST_DIR): | $(BUILD_DIR)
	-mkdir $(DIST_DIR)

$(INT_DIR): | $(BUILD_DIR)
	-mkdir $(INT_DIR)

GOPATH  := $(shell go env GOPATH)
GOCACHE := $(shell go env GOCACHE)
GOBIN   ?= $(GOPATH)/bin

# External binaries
#
# The following external binaries are required
# by this make file.
#
# We will abort any further commands if go
# is not installed.
#
# For docker, docker-compose, etc., we will
# only throw an error when evaluating targets
# that use that functionality and throw
# an error
GOLANGCILINT   := $(GOBIN)/golangci-lint
GOIMPORTS      := $(GOBIN)/goimports
GOCOVMERGE     := $(GOBIN)/gocovmerge
GOCOVXML       := $(GOBIN)/gocov-xml
GOCOV          := $(GOBIN)/gocov
RICHGO         := $(GOBIN)/richgo
MAKEDOC        := $(GOBIN)/makedoc
STATIK         := $(GOBIN)/statik
GORELEASER     := bin/goreleaser/v.0.154.0/$(OS)/goreleaser
GOFUMPT        := $(GOBIN)/gofumpt

$(GOLANGCILINT):
	# To bump, simply change the version at the end to the desired version. The git sha here points to the newest commit
	# of the install script verified by our team located here: https://github.com/golangci/golangci-lint/blob/master/install.sh
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/b90551cdf9c6214075f2a40d1b5595c6b41ffff0/install.sh | sh -s -- -b ${GOBIN} v1.32.2

$(GOIMPORTS):
	$(GO) get -u golang.org/x/tools/cmd/goimports

$(GOCOVMERGE):
	$(GO) get -u github.com/wadey/gocovmerge

$(GOCOVXML):
	$(GO) get -u github.com/AlekSi/gocov-xml

$(GOCOV):
	$(GO) get -u github.com/axw/gocov/gocov

$(RICHGO):
	$(GO) get -u github.com/kyoh86/richgo

$(MAKEDOC):
	$(GO) get -u github.com/paulbes/makedoc

$(STATIK):
	$(GO) get -u github.com/rakyll/statik

$(GOFUMPT):
	$(GO) get -u mvdan.cc/gofumpt

GO := $(shell command -v go 2> /dev/null)
ifndef GO
$(error go is required, please install)
endif

PKGS  = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
FILES = $(shell find . -name '.?*' -prune -o -name vendor -prune -o -name '*.go' -print)

## Release
release-local:
	 docker run --rm --privileged \
	  -v $$PWD:/go/src/github.com/oslokommune/okctl \
	  -v /var/run/docker.sock:/var/run/docker.sock \
	  -w /go/src/github.com/oslokommune/okctl \
	  troian/golang-cross release --config=/go/src/github.com/oslokommune/okctl/.goreleaser-local.yml --snapshot --skip-publish --rm-dist

release:
	 docker run --rm --privileged \
	  -v $$PWD:/go/src/github.com/oslokommune/okctl \
	  -v /var/run/docker.sock:/var/run/docker.sock \
	  -w /go/src/github.com/oslokommune/okctl \
	  -e GITHUB_TOKEN \
	  troian/golang-cross release --rm-dist --release-notes=$(RELEASE_NOTES) --config=/go/src/github.com/oslokommune/okctl/.goreleaser.yml

## Generate
generate: $(STATIK)
	$(GO) generate

## Format
fmt: $(GOFUMPT)
	$(GO) fmt $(PKGS)
	$(GOFUMPT) -s -w $(FILES)

## Imports
imports: $(GOIMPORTS)
	$(foreach gofile,$(FILES),$(GOIMPORTS) -w $(gofile) &&) true

## Linting
lint: $(GOLANGCILINT)
	$(GOLANGCILINT) run

## Testing
TIMEOUT  = 10m
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
            '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
            $(PKGS))
TEST_TARGETS := test-default test-bench test-short test-verbose test-race
test-bench:   ARGS=-run=__absolutelynothing__ -bench=.
test-short:   ARGS=-short
test-verbose: ARGS=-v
test-race:    ARGS=-race
$(TEST_TARGETS): test
check test tests: fmt lint $(RICHGO)
	$(GO) test -timeout $(TIMEOUT) $(ARGS) $(TESTPKGS) | tee >(RICHGO_FORCE_COLOR=1 $(RICHGO) testfilter); \
		test $${PIPESTATUS[0]} -eq 0

test-update:
	$(GO) test ./... -update

integration:
	$(GO) test -tags=integration ./...

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html
test-coverage-tools: | $(GOCOVMERGE) $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(BUILD_DIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint test-coverage-tools
	@mkdir -p $(COVERAGE_DIR)/coverage
	@for pkg in $(TESTPKGS); do \
        go test \
            -coverpkg=$$(go list -f '{{ join .Deps "\n" }}' $$pkg | \
                    grep '^$(MODULE)/' | \
                    tr '\n' ',')$$pkg \
            -covermode=$(COVERAGE_MODE) \
            -coverprofile="$(COVERAGE_DIR)/coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
     done
	@$(GOCOVMERGE) $(COVERAGE_DIR)/coverage/*.cover > $(COVERAGE_PROFILE)
	@$(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@$(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)
