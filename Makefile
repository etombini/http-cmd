
VERSION=$(shell git describe --tags --abbrev=0)
BUILD=$(shell git rev-parse HEAD)

BIN_DIR=$(CURDIR)/bin

HTTPCMD_SRC="cmd/http-cmd/http-cmd.go"
HTTPCMD_BIN="$(BIN_DIR)/http-cmd"

.PHONY: all

all: httpcmd

httpcmd:
	@echo "Building http-cmd Version: $(VERSION) - Build: $(BUILD)"
	mkdir -p $(BIN_DIR)
	go build -ldflags "-X github.com/etombini/http-cmd/pkg/version.version=$(VERSION) -X github.com/etombini/http-cmd/pkg/version.build=$(BUILD)"  -o $(HTTPCMD_BIN) $(HTTPCMD_SRC)
	

clean:
	rm -rf $(BIN_DIR)