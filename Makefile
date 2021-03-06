
VERSION=$(shell git describe --tags --abbrev=0)
BUILD=$(shell git rev-parse HEAD)

BIN_DIR=$(CURDIR)/bin

HTTPCMD_SRC="cmd/http-cmd/http-cmd.go"
HTTPCMD_BIN="$(BIN_DIR)/http-cmd"
HTTPCMD_TEST=$(shell find . -name "*_test.go")

#REPO=github.com/etombini/http-cmd
#GOPATH =$(shell echo "$(CURDIR)" |sed -e 's|\(.*\)/src/$(REPO)|\1|' )

.PHONY: all test $(HTTPCMD_TEST)

all: httpcmd

httpcmd:
	@echo "Building http-cmd Version: $(VERSION) - Build: $(BUILD)"
	@mkdir -p $(BIN_DIR)
	go build -ldflags "-X github.com/etombini/http-cmd/pkg/version.version=$(VERSION) -X github.com/etombini/http-cmd/pkg/version.build=$(BUILD)"  -o $(HTTPCMD_BIN) $(HTTPCMD_SRC)
	

run: httpcmd
	sudo -u http-cmd ./bin/http-cmd -config ./local/http-cmd.yaml

test: $(HTTPCMD_TEST)
	
$(HTTPCMD_TEST):
	go test $@

clean:
	rm -rf $(BIN_DIR)