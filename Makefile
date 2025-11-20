
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=rest-server
TAG=v1.0
WORKDIR=`pwd`
GOBUILDDIR=/go/src/github.com/rest-server
VERSION="v1.0.0"
DBVERSION="v1.0.2"
COMMIT=`git rev-parse --short HEAD`
BUILDDATE=`date -u +%Y-%m-%dT%H:%M:%SZ`
LDFLAGS="-ldflags=\" -X 'github.com/yi-cloud/rest-server/pkg/server.Version=$(VERSION)' -X 'github.com/yi-cloud/rest-server/pkg/server.Commit=$(COMMIT)'	-X 'github.com/yi-cloud/rest-server/pkg/server.BuildDate=$(BUILDDATE)' -X 'github.com/yi-cloud/rest-server/pkg/db.DBVersion=$(DBVERSION)'\""

all: build
build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

docker-all: docker-build docker-image
docker-build:
	docker run --rm -it -e CGO_ENABLED=1 -v $(WORKDIR):$(GOBUILDDIR) -w $(GOBUILDDIR) golang:1.24.4 go build $(LDFLAGS) -o "$(BINARY_NAME)" -v
docker-image:
	docker build -t rest-server:$(TAG) .
