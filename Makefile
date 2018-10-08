.PHONY: all test clean build

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=octopus
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker
TAG = 0.1.1
IMAGE_NAME = octopus
REGISTRY = art-hq.intranet.qualys.com:5001
REPO = datalake

all: deps build
docker: build-image push-image clean-image
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)
deps:
	dep ensure

build-image:
	docker build -t $(REGISTRY)/$(REPO)/$(IMAGE_NAME) -f build/Dockerfile .
	docker tag $(REGISTRY)/$(REPO)/$(IMAGE_NAME) $(REGISTRY)/$(REPO)/$(IMAGE_NAME):$(TAG)

push-image:
	docker push $(REGISTRY)/$(REPO)/$(IMAGE_NAME)
	docker push $(REGISTRY)/$(REPO)/$(IMAGE_NAME):$(TAG)

clean-image:
	docker rmi $(REGISTRY)/$(REPO)/$(IMAGE_NAME):$(TAG) || :
	docker rmi $(REGISTRY)/$(REPO)/$(IMAGE_NAME) || :

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v