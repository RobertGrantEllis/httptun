# This makefile requires:
# * Golang:
# * Dep: 

# Parameters
BINARY_NAME=$(shell basename `pwd`)

# Commands
GOCMD=go
GOFMT=$(GOCMD) fmt
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

DEP=dep
PKILL=pkill

all: format dep test install
format:
	$(GOFMT) ./...
dep:
	$(DEP) ensure
test:
	$(GOTEST) -v ./...
install:
	$(GOINSTALL) -v
clean:
	$(GOCLEAN) -v

start:
	$(BINARY_NAME) serve
stop:
	$(PKILL) $(BINARY_NAME)