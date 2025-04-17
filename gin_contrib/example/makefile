LOCALBIN ?= $(shell pwd)/bin
SWAG ?= $(LOCALBIN)/swag
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

init:
	## 安装 swag 二进制
	GOBIN=$(LOCALBIN) go install github.com/swaggo/swag/cmd/swag@v1.16.3

doc:init
	$(SWAG) fmt
	$(SWAG) init --parseDependency --parseDepth 3