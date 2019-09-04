GOCMD=go
GO11MODULE=on

ARTIFACT_NAME=lagoon
ARTIFACT_DESTINATION=$(GOPATH)/bin

all: test build
build:
	$(GOCMD) build -o $(ARTIFACT_DESTINATION)/$(ARTIFACT_NAME) -v
test:
	$(GOCMD) fmt ./...
	$(GOCMD) vet ./...
	$(GOCMD) test -v ./...
clean:
	$(GOCMD) clean
