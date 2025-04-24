LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLINES ?= $(LOCALBIN)/golines
GOFUMPT ?= $(LOCALBIN)/gofumpt
GOLINTER ?=$(LOCALBIN)/golangci-lint
GOIMPORTS ?=$(LOCALBIN)/goimports-reviser

.PHONY: init
init:
	## 安装 golangci-lint 二进制
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) v2.1.2
	## 安装 gofumpt 二进制
	GOBIN=$(LOCALBIN) go install mvdan.cc/gofumpt@v0.6.0
	## 安装 goimports-reviser 二进制
	GOBIN=$(LOCALBIN) go install github.com/incu6us/goimports-reviser/v3@latest
	## 安装 golines 二进制
	GOBIN=$(LOCALBIN) go install github.com/segmentio/golines@v0.12.2

.PHONY: fmt
fmt:
	$(GOLINES) ./ -m 119 -w --base-formatter gofmt --no-reformat-tags
	$(GOFUMPT) -l -w .
	$(GOIMPORTS) -rm-unused -set-alias -format ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: fmt vet
	$(GOLINTER) run

.PHONY: test
test:
	go test ./...


