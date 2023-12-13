APPNAME := $(notdir $(CURDIR))

GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
GIT_SHA := $(shell git rev-parse --short HEAD)
VERSION := $(GIT_TAG)-$(GIT_SHA)$(shell git diff --quiet || echo "-dirty")

CMDPATH := ./cmd/$(APPNAME)
BUILDPATH := ./build

# Go parameters
GOVERSION=1.21.1
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get

LDFLAGS := -ldflags '-s -w \
						-X=github.com/mikejoh/$(APPNAME)/internal/version.Version=$(VERSION) \
						-X=github.com/mikejoh/$(APPNAME)/internal/version.Name=$(APPNAME)'

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -v -o $(BUILDPATH)/$(APPNAME) $(CMDPATH)

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(BUILDPATH)/$(APPNAME)

run:
	$(GOBUILD) -v -o $(BUILDPATH)/$(APPNAME) $(CMDPATH)
	$(BUILDPATH)/$(APPNAME)

install:
	cp $(BUILDPATH)/$(APPNAME) ~/.local/bin

deps:
	$(GOGET) ./...

.PHONY: all build