# Makefile for rivermq

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
BINARY=rivermq
VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

.DEFAULT_GOAL := integration


.PHONY: install
install:
	go install $(SOURCEDIR)/...


.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f .coverage.txt ] ; then rm .coverage.txt ; fi
	find ./ -name .coverprofile -print0 | xargs -0 rm


test: $(SOURCES)
	go test ${SOURCEDIR}
	go test ${SOURCEDIR}/model


.PHONY: integration
integration:
	go test .         -coverprofile=.coverprofile           -covermode=atomic
	go test ./model   -coverprofile=./model/.coverprofile   -covermode=atomic -tags integration
	go test ./route   -coverprofile=./route/.coverprofile   -covermode=atomic -tags integration
	go test ./handler -coverprofile=./handler/.coverprofile -covermode=atomic -tags integration
	go test ./deliver -coverprofile=./deliver/.coverprofile -covermode=atomic -tags integration
	gover . .coverage.txt


build: $(SOURCES)
	go build -o ${BINARY} $(SOURCEDIR)/*.go


run:
	go build -o ${BINARY} $(SOURCEDIR)/*.go
	$(SOURCEDIR)/$(BINARY)
