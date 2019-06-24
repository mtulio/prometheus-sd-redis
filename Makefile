
# Globals

DIR_BIN ?= ./bin
DIR_DIST ?= ./dist
DIR_TMP ?= ./.build

BIN_NAME ?= prometheus-sd-redis

REPO_OWNER ?= mtulio
REPO_NAME ?= prometheus-sd-redis

## Build flags
export CGO_ENABLED=0

BUILD_TAG ?= "latest"
BUILD_COMMIT_SHA ?= $(shell git rev-parse HEAD)
BUILD_OS ?= linux
BUILD_ARCH ?= amd64
BUILD_DIR ?= $(DIR_TMP)/$(BIN_NAME)-$(BUILD_OS)-$(BUILD_ARCH)
BUILD_BIN ?= $(DIR_TMP)/$(BIN_NAME)-$(BUILD_OS)-$(BUILD_ARCH)/$(BIN_NAME)

GO_LDFLAGS ?= -s -w -extldflags \"-static\" -X main.BuildVersion=$(BUILD_TAG) -X main.BuildCommitSha=$(BUILD_COMMIT_SHA) -X main.BuildDate=$(shell date +%F-%T)

GOX_OS ?= "darwin linux windows freebsd netbsd openbsd"
GOX_ARCH ?= "386 amd64"

## Modules
export GO111MODULE=on


# Dependencies

install-ghr:
	@echo installing ghr
	GO111MODULE=0 go get github.com/tcnksm/ghr

install-gox:
	@echo installing gox
	GO111MODULE=0 go get github.com/mitchellh/gox

ensure-dirs:
	@mkdir -p $(BUILD_DIR)
	@mkdir -p $(DIR_TMP)/$(BIN_NAME)-$(BUILD_OS)-$(BUILD_ARCH)

modules-dep:
	if [[ -f 'go.mod' ]]; then \
		go mod tidy; \
	fi

dependencies: install-ghr install-gox

# Build

fix-fmt:
	go fmt

check-bin-amd64:
	@if [ -f $(DIR_BIN)/$(BIN_NAME) ]; then \
		echo -n "\n\tSUCCESS, binary available on $(DIR_BIN)/$(BIN_NAME)\n"; \
	fi
	@chmod +x $(DIR_BIN)/$(BIN_NAME)

check-bin-cross:
	@echo TODO check cross bins

build-amd64: ensure-dirs clear-bins
	go build \
		-ldflags $(GO_LDFLAGS) \
		-o $(BUILD_BIN) && \
		mv $(BUILD_BIN) $(DIR_BIN)/

build: build-amd64 check-bin-amd64

build-cross-all: ensure-dirs clear-bins
	gox -verbose -rebuild \
		-os=$(GOX_OS) \
		-arch=$(GOX_ARCH) \
		-ldflags "${GO_LDFLAGS}" \
		-output "$(DIR_TMP)/$(BIN_NAME)-${BUILD_TAG}.{{.OS}}-{{.Arch}}/{{.Dir}}"

build-cross: build-cross-all check-bin-cross

# Publish

define binPack
	echo packBin: $1; \
	cp LICENSE README.md "$(DIR_TMP)/$1/"; \
	tar -C "$(DIR_TMP)/" -czf "$(DIR_DIST)/$1.tar.gz" "$1"
endef

define binSum
	cd $(DIR_DIST); \
	sha256sum *.gz *.zip > sha256sums.txt; \
	cd -
endef

define binPublish
	ghr -u $(REPO_OWNER) \
		-r $(REPO_NAME) \
		--replace "${BUILD_TAG}" \
		dist/
endef

publish-github:
	@echo "Creating  dir $(DIR_DIST)"
	mkdir -p $(DIR_DIST) || true

	$(foreach var,$(shell ls $(DIR_TMP)),$(call binPack,$(var));)
	$(call binSum)
	$(call binPublish)

# Cleaner

clear-bins:
	@rm -f $(DIR_BIN)/$(BIN_NAME)
	@rm -fr $(DIR_TMP)
	@rm -fr $(DIR_DIST)

clean: clear-bins
