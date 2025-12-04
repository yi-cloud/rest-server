
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=auth-server
TAG ?= v1.0
WORKDIR=`pwd`
GOBUILDDIR=/go/src/github.com/rest-server
CMDDIR=cmd/auth
VERSION ?= v1.0.0
DBVERSION ?= v1.0.2
CLUSTERID ?= 00000000-0000-0000-0000-000000000000
PRODUCT ?= auth
COMMIT=`git rev-parse --short HEAD`
BUILDDATE=`date -u +%Y-%m-%dT%H:%M:%SZ`
PKG=github.com/yi-cloud/rest-server/pkg
LDFLAGS=-ldflags=" -X '$(PKG)/server.Version=$(VERSION)' -X '$(PKG)/server.Commit=$(COMMIT)' -X '$(PKG)/server.BuildDate=$(BUILDDATE)' -X '$(PKG)/db.DBVersion=$(DBVERSION)' -X '$(PKG)/license.ClusterId=$(CLUSTERID)' -X '$(PKG)/license.Product=$(PRODUCT)'"

all: build
build:
	cd $(GOBUILDDIR)/$(CMDDIR)
	GO111MODULE=on GOPROXY=https://goproxy.cn,direct CGO_ENABLED=1 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
	mv $(CMDDIR)/$(BINARY_NAME) .
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)

docker-all: docker-build docker-image
docker-build:
	docker run --rm -it -e CGO_ENABLED=1 -e GO111MODULE=on -e GOPROXY=https://goproxy.cn,direct -v $(WORKDIR):$(GOBUILDDIR) -w $(GOBUILDDIR)/$(CMDDIR) golang:1.24.4 go build $(LDFLAGS) -o "$(BINARY_NAME)" -v
	mv $(CMDDIR)/$(BINARY_NAME) .
docker-image:
	docker build -t auth-server:$(TAG) .
