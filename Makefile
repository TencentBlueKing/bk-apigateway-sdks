LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLINTER ?=$(LOCALBIN)/golangci-lint

.PHONY: init
init:
	## 安装 golangci-lint 二进制
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v2.14.0/install.sh | sh -s -- -b $(LOCALBIN) v2.14.0

.PHONY: fmt
fmt:
	$(GOLINTER) fmt

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: fmt vet
	$(GOLINTER) run

.PHONY: test
test:
	go test ./...


