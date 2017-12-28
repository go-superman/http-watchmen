NAME = http-watchmen
BINARY = http-watchmen

GO_FLAGS = #-v -race
GO_VERSION = latest
APHLIE_PACK = -tags netgo

GOOS = `go env GOHOSTOS`
GOARCH = `go env GOHOSTARCH`

SOURCE_DIR = ./

all: build

.PHONY : clean build linux fmt docker

clean:
	go clean -i $(GO_FLAGS) $(SOURCE_DIR)
	rm -f $(BINARY)
	rm -rf linux

fmt:
	goimports -w ...

build:
	mkdir -p build/$(GOOS)-$(GOARCH)
	go build $(GO_LDFLAGS) $(GO_FLAGS) $(APHLIE_PACK) -o build/$(GOOS)-$(GOARCH)/$(BINARY) $(SOURCE_DIR)

package:
	cd build/$(GOOS)-$(GOARCH)/ &&  tar zcvf $(NAME)-$(GOOS)-$(GOARCH)-`git describe --tags`.tar.gz $(BINARY)

linux:
	mkdir -p build/linux-amd64
	GOOS=linux GOARCH=amd64 go build $(GO_LDFLAGS) $(GO_FLAGS)  $(APHLIE_PACK) -o linux/$(BINARY) $(SOURCE_DIR)

docker:
	docker run --rm -v "`pwd`":/go/src/$(NAME) -w /go/src/$(NAME) golang:$(GO_VERSION) bash -c "make build "
