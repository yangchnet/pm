$(if $(shell ! which git &>/dev/null),$(error please install git first))

# Set a specific PALTFORM
ifeq ($(origin PLATFORM),undefined)
	ifeq ($(origin GOOS),undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH),undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM  := $(GOOS)_$(GOARCH)
	# Use linux as the default OS when building images
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN),undefined)
	GOBIN := $(GOPATH)/bin
endif

## go.build.linux_amd64.<service>
.PHONY: build.%
build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "==========> Building binary $(COMMAND) for $(GOOS) $(GOARCH)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$(COMMAND)$(GO_OUT_EXT) main.go

## build: 编译所有服务为二进制可执行文件
.PHONY: build
build: $(addprefix build., $(addprefix $(PLATFORM)., pm))

ALL: build
