-include .env
-include .env.local
# ref: https://kodfabrik.com/journal/a-good-makefile-for-go/

PROJECTNAME       = $(shell basename "$(PWD)")
DEFAULT_DOC_NAME ?= doc.go
APPNAME           = $(shell grep -E "appName[ \t]+=[ \t]+" $(DEFAULT_DOC_NAME)|grep -Eo "\\\".+\\\"")
VERSION           = $(shell grep -E "version[ \t]+=[ \t]+" $(DEFAULT_DOC_NAME)|grep -Eo "[0-9.]+")

#GIT_VERSION  := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
#GIT_REVISION := $(shell git rev-parse --short HEAD)
#GIT_HASH     =  $(shell git rev-parse HEAD)
#BUILDTIME   := $(shell date "+%Y%m%d_%H%M%S")
#BUILDTIME   =  $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
#BUILDTIME    =  $(shell date -u '+%Y-%m-%d_%H-%M-%S')
#GOVERSION    =  $(shell go version)

# Go related variables.
GOBASE = $(shell pwd)
#GOPATH="$(GOBASE)/vendor:$(GOBASE)"
#GOPATH=$(GOBASE)/vendor:$(GOBASE):$(shell dirname $(GOBASE))
GOPATH2= $(shell dirname $(shell dirname $(shell dirname $(GOBASE))))
GOPATH1= $(shell dirname $(GOPATH2))
ifneq ($(wildcard $(GOPATH2)/src),)
	GOPATH  ?= $(GOPATH2)
else
	GOPATH  ?= $(HOME)/go
endif
GOBIN       ?= $(GOBASE)/bin
GOFILES     ?= $(wildcard *.go)
BIN         ?= $(GOPATH)/bin
GOLINT      ?= $(BIN)/golint
GOCYCLO     ?= $(BIN)/gocyclo
GOYOLO      ?= $(BIN)/yolo

GO111MODULE ?= $(or $(shell go env GO111MODULE),on)
GOPROXY     ?= $(or $(GOPROXY_CUSTOM),$(or "$(shell go env GOPROXY)","https://goproxy.io,direct"))
# GOPROXY=https://goproxy.io # https://athens.azurefd.net

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID=/tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent


M = $(shell printf "\033[34;1m▶\033[0m")
CN = hedzr/$(N)
ADDR = ":5q5q"
SERVER_START_ARG=server run
SERVER_STOP_ARG=server stop


goarch=amd64
W_PKG=github.com/hedzr/cmdr/conf
TIMESTAMP=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
TIMESTAMP=$(shell date -u '+%Y-%mm-%ddT%HH:%MM:%SS')
GOVERSION=$(shell go version)
GIT_VERSION  := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
GIT_REVISION := $(shell git rev-parse --short HEAD)
GIT_SUMMARY=$(shell git describe --tags --dirty --always)
GIT_DESC=$(shell git log --oneline -1)
GIT_HASH=$(shell git rev-parse HEAD)
GOBUILD_TAGS ?= "-tags='hzstudio sec antonal'"
LDFLAGS=-s -w -X '$(W_PKG).Buildstamp=$(TIMESTAMP)' -X '$(W_PKG).GIT_HASH=$(GIT_REVISION)' -X '$(W_PKG).GitSummary=$(GIT_SUMMARY)' -X '$(W_PKG).GitDesc=$(GIT_DESC)' -X '$(W_PKG).BuilderComments=$(BUILDER_COMMENT)' -X '$(W_PKG).GoVersion=$(GOVERSION)' -X '$(W_PKG).Version=$(VERSION)' -X '$(W_PKG).AppName=$(APPNAME)'



GRPC_GATEWAY  = $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
GO_ANNOTATION = google/api/annotations.proto


export PATH := $(BIN):$(PATH)




goarch=$(shell go env GOARCH)
goos=$(shell go env GOOS)
defgoarch=$(shell go env GOARCH)
defgoos=$(shell go env GOOS)
W_PKG=github.com/hedzr/cmdr/conf
LDFLAGS := -s -w \
	-X '$(W_PKG).Buildstamp=$(BUILDTIME)' \
	-X '$(W_PKG).GIT_HASH=$(GIT_REVISION)' \
	-X '$(W_PKG).GoVersion=$(GOVERSION)' \
	-X '$(W_PKG).Version=$(VERSION)'
# -X '$(W_PKG).AppName=$(APPNAME)'
GOSYS := GOARCH="$(goarch)" GOOS="$(os)" \
	GOPATH="$(GOPATH)" GOBIN="$(BIN)" \
	GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) CGO_ENABLED=$(CGO_ENABLED) go
CGO := GOARCH="$(goarch)" GOOS="$(os)" \
	GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" \
	GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) CGO_ENABLED=1 go
GO := GOARCH="$(goarch)" GOOS="$(os)" \
	GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" \
	GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) CGO_ENABLED=$(CGO_ENABLED) go
GO_OFF := GOARCH="$(goarch)" GOOS="$(os)" \
	GOPATH="$(GOPATH)" GOBIN="$(GOBIN)" \
	GO111MODULE=off CGO_ENABLED=$(CGO_ENABLED) go


ifeq ($(OS),Windows_NT)
    LS_OPT=
    CCFLAGS += -D WIN32
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        CCFLAGS += -D AMD64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            CCFLAGS += -D AMD64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            CCFLAGS += -D IA32
        endif
    endif
else
    LS_OPT=
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        OS = Linux
        CCFLAGS += -D LINUX
        LS_OPT=--color
    endif
    ifeq ($(UNAME_S),Darwin)
        OS = macOS
        CCFLAGS += -D OSX
        LS_OPT=-G
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif


ifeq ($(OS),macOS)
	TIMESTAMP=$(shell date -Iseconds)
endif




.PHONY: bgo swagger pb
.PHONY: tidy clean clean2 tools directories
.PHONY: build compile exec run
.PHONY: build-windows build-win build-linux build-nacl build-plan9 build-freebsd build-darwin build-m1
.PHONY: build-ci go-build go-generate go-mod-download go-get go-install go-clean
.PHONY: docker godoc format fmt lint cov gocov coverage codecov cyclo bench


## bgo: compile proto buffer
bgo:
	@echo "Running bgo building ..."
	bgo -s

## swagger: rebuild swagger doc files
swagger:
	swag init --output ./cli/atonal/swaggerdocs

## pb: compile proto buffer
pb: | $(BASE) $(GRPC_GATEWAY) $(BIN)/protoc-gen-go $(BIN)/protoc-gen-swagger $(BIN)/protoc-gen-grpc-gateway $(GO_ANNOTATION)
	INC="-I/usr/local/include -I. -I./v1 \
		 -I$(GOPATH)/src \
		 -I$(GRPC_GATEWAY)"
	protoc $(INC) --go_out=plugins=grpc:. v1/*.proto
	protoc $(INC) --grpc-gateway_out=logtostderr=true:. v1/*.proto
	protoc $(INC) --swagger_out=logtostderr=true:. v1/*.proto

$(GRPC_GATEWAY): | $(BASE)
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go get -v github.com/grpc-ecosystem/grpc-gateway
	$(MAKE) tidy

google/api:
	@-mkdir -p google/api

$(GO_ANNOTATION): | google/api
	curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > google/api/annotations.proto
	curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > google/api/http.proto

$(BIN)/protoc-gen-grpc-gateway:
	go install -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest

$(BIN)/protoc-gen-swagger:
	go install -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

$(BIN)/protoc-gen-go:
	go install -v github.com/golang/protobuf/protoc-gen-go@latest


## tidy: Go Module Tidy
tidy:
	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go mod tidy


## clean: Clean build files. Runs `go clean` internally.
clean:
	# @(MAKEFILE) go-clean
	@-rm v1/*.pb.go v1/*.pb.gw.go v1/*.json 2>/dev/null

## clean2: Clean build files. Runs `go clean` internally.
clean2:
	@(MAKEFILE) go-clean



#---------------------


## build: Compile the binary. Synonym of `compile`
build: compile

## compile: Compile the binary.
compile: directories go-clean go-generate
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-build 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/' 1>&2

go-build:
	@echo "     With os=$(goos) arch=$(goarch)"
	@-$(MAKE) -s go-build-task os="$(goos)" goarchset="$(goarch)"


## build-win: or build-windows, build to windows executable, for LAN deploy manually.
build-windows: build-win
build-win:
	@-$(MAKE) -s go-build-task os=windows goarchset=amd64

## build-linux: build to linux executable, for LAN deploy manually.
build-linux:
	@-$(MAKE) -s go-build-task os=linux goarchset=amd64

## build-nacl: build to nacl executable, for LAN deploy manually.
build-nacl:
	# NOTE: can't build to nacl with golang 1.14 and darwin
	#    chmod +x $(GOBIN)/$(an)_$(os)_$(goarch)*;
	#    ls -la $(LS_OPT) $(GOBIN)/$(an)_$(os)_$(goarch)*;
	#    gzip -f $(GOBIN)/$(an)_$(os)_$(goarch);
	#    ls -la $(LS_OPT) $(GOBIN)/$(an)_$(os)_$(goarch)*;
	@-$(MAKE) -s go-build-task os=nacl goarchset="386 arm amd64p32"
	@echo "  < All Done."
	@ls -la $(LS_OPT) $(GOBIN)/*

## build-plan9: build to plan9 executable, for LAN deploy manually.
build-plan9: goarchset = "386 amd64"
build-plan9:
	@-$(MAKE) -s go-build-task os=plan9 goarchset=$(goarchset)

## build-freebsd: build to freebsd executable, for LAN deploy manually.
build-freebsd:
	@-$(MAKE) -s go-build-task os=freebsd goarchset=amd64

## build-riscv: build to riscv64 executable, for LAN deploy manually.
build-riscv:
	@-$(MAKE) -s go-build-task os=linux goarchset=riscv64

## build-ci: run build-ci task. just for CI tools
build-ci:
	@echo "  >  Building binaries in CI flow..."
	$(foreach os, linux darwin windows, \
	  @-$(MAKE) -s go-build-task os=$(os) goarchset="386 amd64" \
	)
	@-$(MAKE) -s go-build-task os="darwin" goarchset="arm64"
	@echo "  < All Done."
	@ls -la $(LS_OPT) $(GOBIN)/*

## build-darwin: build to riscv64 executable, for LAN deploy manually.
build-darwin:
	@-$(MAKE) -s go-build-task os=darwin goarchset=amd64

## build-m1: run build-ci task. just for CI tools
build-m1:
	@-$(MAKE) -s go-build-task os="darwin" goarchset="arm64"


go-build-task: directories tools go-generate
	@echo "-->--"
	@echo "  >  Building $(os)/$(goarchset) binary..."
	# unsupported GOOS/GOARCH pair nacl/386 ??
	$(foreach an, $(MAIN_APPS), \
	  $(foreach san, $(SUB_APPS), \
		  $(eval DOCNAME = "./$(an)/$(san)/$(DEFAULT_SUBAPP_DOC_NAME).go") \
		  $(eval MAINGONAME = "./$(an)/$(san)/") \
		  if [ -f "$(DOCNAME)" ]; then \
	        $(eval APPNAME = $(patsubst "%",%,$(shell grep -E "appName[ \t]+=[ \t]+" "$(DOCNAME)" 2>/dev/null|grep -Eo "\\\".+\\\""))) \
	        $(eval VERSION = $(shell grep -E "version[ \t]+=[ \t]+" "$(DOCNAME)" 2>/dev/null|grep -Eo "[0-9.]+")) \
	        $(foreach goos, $(os), \
	        $(foreach goarch, $(goarchset), \
		    echo "     > DOCNAME = $(DOCNAME), MAINGONAME = $(MAINGONAME)"; \
		    echo "     > APPNAMEs = appname:$(APPNAME)|projname:$(PROJECTNAME)|an:$(an)"; \
		    $(MAKE) -s go-build-child os=$(goos) goarch=$(goarch) SUFFIX="_$(goos)-$(goarch)" DOCNAME=$(DOCNAME) MAINGONAME=$(MAINGONAME) an=$(an) san=$(san); \
		    ) ) \
		  else \
	        $(eval APPNAME = $(san)) \
	        $(eval VERSION = $(shell grep -iE "version[ \t]+=[ \t]+" "$(DEFAULT_DOC_NAME)" 2>/dev/null|grep -Eo "[0-9.]+")) \
	        $(foreach goos, $(os), \
	        $(foreach goarch, $(goarchset), \
	        echo "        APPNAME=$(san), VERSION=$(VERSION)"; \
		    $(MAKE) -s go-build-child os="$(goos)" goarch="$(goarch)" \
		    SUFFIX="_$(goos)-$(goarch)" DOCNAME="$(DEFAULT_DOC_NAME)" \
		    MAINGONAME=$(MAINGONAME) an=$(an) san=$(san) APPNAME=$(APPNAME) VERSION=$(VERSION); \
		    ) ) \
		  fi; \
	  ) \
	)
	#	$(foreach an, $(MAIN_APPS), \
	#	  $(eval ANAME := $(shell if [ "$(an)" == "cli" ]; then echo $(APPNAME); else echo $(an); fi; )) \
	#	  echo "  >  APP NAMEs = appname:$(APPNAME)|projname:$(PROJECTNAME)|an:$(an)|ANAME:$(ANAME)"; \
	#	  $(foreach goarch, $(goarchset), \
	#	    echo "     >> Building (-trimpath) $(GOBIN)/$(ANAME)_$(os)_$(goarch)...$(os)" >/dev/null; \
	#	    $(GO) build -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(ANAME)_$(os)_$(goarch) $(GOBASE)/$(MAIN_BUILD_PKG)/$(an); \
	#	    chmod +x $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	#	    ls -la $(LS_OPT) $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	#	) \
	#	)
	#@ls -la $(LS_OPT) $(GOBIN)/*linux*

go-build-child:
	@echo "     >  go-build-child: suffix = $(SUFFIX), DOCNAME = $(DOCNAME), MAINGONAME = $(MAINGONAME), AN = $(an), san = $(san)."
	@echo "        > detecting SUB_APPS: $(APPNAME) v$(VERSION), def: ($(defgoos), $(defgoarch)), curr: ($(os), $(goarch))"
	$(eval ANAME = $(shell for an1 in $(MAIN_APPS); do \
	    if [[ "$(an)" == $$an1 ]]; then \
	      if [[ $$an1 == cli ]]; then \
	        if [[ $(san) == "$(APPNAME)" ]]; then echo $(APPNAME); else echo $(san)$(SUFFIX); fi; \
	      else \
	          if [ "$(os)" = "$(defgoos)" -a "$(goarch)" = "$(defgoarch)" ]; then \
	            echo $(san); \
	          else \
	            echo $(san)$(SUFFIX); \
	          fi; \
	      fi; \
	    fi; \
	done))
	$(eval LDFLAGS = -s -w \
	    -X '$(W_PKG).Buildstamp=$(TIMESTAMP)' \
	    -X '$(W_PKG).GIT_HASH=$(GIT_REVISION)' \
	    -X '$(W_PKG).GitSummary=$(GIT_SUMMARY)' \
	    -X '$(W_PKG).GitDesc=$(GIT_DESC)' \
	    -X '$(W_PKG).BuilderComments=$(BUILDER_COMMENT)' \
	    -X '$(W_PKG).GoVersion=$(GOVERSION)' \
	    -X '$(W_PKG).Version=$(VERSION)' )
	echo "     >  >  Building $(MAINGONAME) -> $(san)$(SUFFIX) v$(VERSION) ..."
	echo "           +race. -trimpath -gcflags=all='-l -B'. appName = $(ANAME), "
	echo "           LDFLAGS = $(LDFLAGS)"
	echo "           TAGS = $(GOBUILD_TAGS)"
	$(GO) build -trimpath -gcflags=all='-l -B' $(GOBUILD_TAGS) -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(ANAME) $(MAINGONAME)
	ls -la $(LS_OPT) $(GOBIN)/$(ANAME)
	@echo "  > go-build-child: END."



# @cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2
#@if [[ -z "$(STDERR)" ]]; then echo; else echo -e "\n\nError:\n\n"; cat $(STDERR)  1>&2; fi

# ## exec: Run given cmd, wrapped with custom GOPATH. eg; make exec run="go test ./..."
#exec:
#	@GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
#	$(run)

#ooo_test:
#	$(eval ANAME := $(shell for an in $(MAIN_APPS); do \
#		if [[ $$an == cli ]]; then A=$(APPNAME); echo $(APPNAME); \
#		else A=$$an; echo $$an; \
#		fi; \
#	done))
#	@echo "ANAME = $(ANAME), $$ANAME, $$A"
#
#oop_test: go-clean go-generate
#	$(MAKE) -s go-build
#
# ## run: go run xxx
#run:
#	@$(GO) run -ldflags "$(LDFLAGS)" $(GOBASE)/cli/main.go


go-build-task-another-one: directories go-generate
	@echo "  >  Building $(os)/$(goarchset) binary..."
	@#echo "  >  LDFLAGS = $(LDFLAGS)"
	# unsupported GOOS/GOARCH pair nacl/386 ??
	$(foreach an, $(MAIN_APPS), \
	  echo "  >  APP NAMEs = appname:$(APPNAME)|projname:$(PROJECTNAME)|an:$(an)"; \
		$(eval ANAME := $(shell for an1 in $(MAIN_APPS); do \
			if [[ $(an) == $$an1 ]]; then \
			  if [[ $$an1 == cli ]]; then echo $(APPNAME); else echo $$an1; fi; \
			fi; \
		done)) \
	  $(foreach goarch, $(goarchset), \
	    echo "     >> Building (-trimpath) $(GOBIN)/$(ANAME)_$(os)_$(goarch)...$(os)" >/dev/null; \
	    pushd "$(MAIN_BUILD_PKG)/$(an)"; \
        $(GO) build -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(ANAME)_$(os)_$(goarch) $(GOBASE)/$(MAIN_BUILD_PKG)/$(an)/$(MAIN_ENTRY_FILE); \
	    popd; chmod +x $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	    ls -la $(LS_OPT) $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	) \
	)
	#	$(foreach an, $(MAIN_APPS), \
	#	  $(eval ANAME := $(shell if [ "$(an)" == "cli" ]; then echo $(APPNAME); else echo $(an); fi; )) \
	#	  echo "  >  APP NAMEs = appname:$(APPNAME)|projname:$(PROJECTNAME)|an:$(an)|ANAME:$(ANAME)"; \
	#	  $(foreach goarch, $(goarchset), \
	#	    echo "     >> Building (-trimpath) $(GOBIN)/$(ANAME)_$(os)_$(goarch)...$(os)" >/dev/null; \
	#	    $(GO) build -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(ANAME)_$(os)_$(goarch) $(GOBASE)/$(MAIN_BUILD_PKG)/$(an); \
	#	    chmod +x $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	#	    ls -la $(LS_OPT) $(GOBIN)/$(ANAME)_$(os)_$(goarch)*; \
	#	) \
	#	)
	#@ls -la $(LS_OPT) $(GOBIN)/*linux*

go-build-1: # never used
	@echo "  >  Building apps: $(MAIN_APPS)..."
	$(foreach an, $(MAIN_APPS), \
		$(eval ANAME := $(shell for an1 in $(MAIN_APPS); do \
			if [[ "$(an)" == $$an1 ]]; then \
			  if [[ $$an1 == cli ]]; then echo $(APPNAME); else echo $$an1; fi; \
			fi; \
		done)) \
	  echo "  >  >  Building $(MAIN_BUILD_PKG)/$(an) -> $(ANAME) ..."; \
	  echo "        +race. -trimpath. APPNAME = $(APPNAME), LDFLAGS = $(LDFLAGS)"; \
	  pushd "$(MAIN_BUILD_PKG)/$(an)"; \
	  $(GO) build -v -race -ldflags "$(LDFLAGS)" -o $(GOBIN)/$(ANAME) $(GOBASE)/$(MAIN_BUILD_PKG)/$(an)/$(MAIN_ENTRY_FILE); \
	  popd; ls -la $(LS_OPT) $(GOBIN)/$(ANAME); \
	)
	ls -la $(LS_OPT) $(GOBIN)/
	if [[ -d ./plugin/demo ]]; then \
	  $(GO) build -v -race -buildmode=plugin -o ./ci/local/share/fluent/addons/demo.so ./plugin/demo && \
	  chmod +x ./ci/local/share/fluent/addons/demo.so && \
	  ls -la $(LS_OPT) ./ci/local/share/fluent/addons/demo.so; fi
	# go build -o $(GOBIN)/$(APPNAME) $(GOFILES)
	# chmod +x $(GOBIN)/*

go-generate:
	@echo "  >  Generating dependency files ('$(generate)') ..."
	@$(GO) generate $(generate) ./...

go-mod-download:
	@$(GO) mod download

go-get:
	# Runs `go get` internally. e.g; make install get=github.com/foo/bar
	@echo "  >  Checking if there is any missing dependencies...$(get)"
	@$(GO) get $(get)

go-install:
	@$(GO) install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	go clean


## docker: docker build
docker:
	@echo "  >  Docker building ..."
	@if [ -n "$(DOCKER_APP_NAME)" ]; then \
	  $(MAKE) -s docker-one APPNAME=$(DOCKER_APP_NAME) VERSION=$(VERSION); \
	else \
	  if [ -n "$(DOCKER_APP_NAMES)" ]; then \
	    $(foreach an, $(DOCKER_APP_NAMES), \
	      $(MAKE) -s docker-one APPNAME=$(an) VERSION=$(VERSION); \
	    ) \
	  else \
	    echo "  >  docker build not available since DOCKER_APP_NAME is empty"; \
	  fi; \
	fi

docker-one:
	echo "  >  docker build $(APPNAME):$(VERSION)..."
	$(eval MAKEFLAGS -= --silent)
	$(eval DOCKER_ORG_NAME ?= hedzr)
	@#   --network=host
	docker build -f app.Dockerfile \
	     --build-arg CN=1 \
	     --build-arg PORT="" \
	     --build-arg APPNAME="$(APPNAME)" \
	     --build-arg VERSION="$(VERSION)" \
	     --build-arg TIMESTAMP="$(TIMESTAMP)" \
	     --build-arg GIT_REVISION="$(GIT_REVISION)" \
	     --build-arg GIT_SUMMARY="$(GIT_SUMMARY)" \
	     --build-arg GIT_DESC="$(GIT_DESC)" \
	     --build-arg BUILDER_COMMENT="$(BUILDER_COMMENT)" \
	     -t $(DOCKER_ORG_NAME)/$(APPNAME):latest \
		 -t $(DOCKER_ORG_NAME)/$(APPNAME):$(VERSION) \
	     .


.PHONY: tools

tools: $(BIN)/golangci-lint

$(BIN)/govulncheck:
	@echo "  >  installing govulncheck ..."
	@$(GO) install -v golang.org/x/vuln/cmd/govulncheck@latest
	[ -x ./bin/govulncheck ] && mv ./bin/govulncheck $(BIN)/

$(BIN)/golangci-lint: | $(GOBASE)
	@echo "  >  installing golangci-lint ..."
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	[ -x ./bin/golangci-lint ] && mv ./bin/golangci-lint $(BIN)/

$(BIN)/golint: | $(GOBASE)   # # # ❶
	@echo "  >  installing golint ..."
	#@-mkdir -p $(GOPATH)/src/golang.org/x/lint/golint
	#@cd $(GOPATH)/src/golang.org/x/lint/golint
	#@pwd
	#@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	#go get -v golang.org/x/lint/golint
	@echo "  >  installing golint ..."
	@$(GO) install golang.org/x/lint/golint
	@cd $(GOBASE)

$(BIN)/gocyclo: | $(GOBASE)  # # # ❶
	@echo "  >  installing gocyclo ..."
	@$(GO) install github.com/fzipp/gocyclo

$(BIN)/yolo: | $(GOBASE)     # # # ❶
	@echo "  >  installing yolo ..."
	@$(GO) install github.com/azer/yolo

$(BIN)/godoc: | $(GOBASE)     # # # ❶
	@echo "  >  installing godoc ..."
	@$(GO) install golang.org/x/tools/cmd/godoc

$(BIN)/gofumpt: | $(GOBASE)
	@echo "  >  installing gofumpt ..."
	@$(GO) install mvdan.cc/gofumpt@latest
	[ -x ./bin/gofumpt ] && mv ./bin/gofumpt $(BIN)/

$(BASE):
	# @mkdir -p $(dir $@)
	# @ln -sf $(CURDIR) $@


## godoc: run godoc server at "localhost;6060"
godoc: | $(GOBASE) $(BIN)/godoc
	@echo "  >  PWD = $(shell pwd)"
	@echo "  >  started godoc server at :6060: http://localhost:6060/pkg/github.com/hedzr/$(PROJECTNAME1) ..."
	@echo "  $  cd $(GOPATH_) godoc -http=:6060 -index -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps"
	( cd $(GOPATH_) && pwd && godoc -v -index -http=:6060 -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps -goroot .; )
	# https://medium.com/@elliotchance/godoc-tips-tricks-cda6571549b


## godoc1: run godoc server at "localhost;6060"
godoc1: # | $(GOBASE) $(BIN)/godoc
	@echo "  >  PWD = $(shell pwd)"
	@echo "  >  started godoc server at :6060: http://localhost:6060/pkg/github.com/hedzr/$(PROJECTNAME1) ..."
	#@echo "  $  GOPATH=$(GOPATH) godoc -http=:6060 -index -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps"
	godoc -v -index -http=:6060 -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps # -goroot $(GOPATH)
	# gopkg.in/hedzr/errors.v2.New
	# -goroot $(GOPATH) -index
	# https://medium.com/@elliotchance/godoc-tips-tricks-cda6571549b

## fmt: =`format`, run gofumpt/gofmt tool
fmt: format
format: | $(GOBASE) $(BIN)/gofumpt
	@echo "  >  gofmt ..."
	@# GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	# gofmt -l -w -s .
	$(BIN)/gofumpt -l -w .

## lint: run golangci-lint/golint tool
#
# https://github.com/golangci/golangci-lint/blob/master/.golangci.reference.yml
# https://golangci-lint.run/
#
# https://freshman.tech/linting-golang/
#
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
lint: | $(GOBASE) # $(GOLINT)
	@echo "  >  golint/golangci-lint ..."
	@# GOPATH=$(GOPATH) GOBIN=$(BIN) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	# $(GOLINT) ./...
	@$(BIN)/golangci-lint run

## vuln: =`vuln-check`, run govulncheck
vuln: vuln-check
vuln-check: | $(BIN)/govulncheck
	$(BIN)/govulncheck ./...

## cov: =`coverage`, run go coverage test
cov: coverage

## gocov: =`coverage`, run go coverage test
gocov: coverage

# coverage: run go coverage test
coverage: | $(GOBASE)
	@echo "  >  gocov ..."
	@echo $(CGO) test $(COVER_TEST_TARGETS) -v -race -coverprofile=coverage.txt -covermode=atomic -timeout=20m -test.short | tee coverage.log
	@$(CGO) test $(COVER_TEST_TARGETS) -v -race -coverprofile=coverage.txt -covermode=atomic -timeout=20m -test.short | tee coverage.log
	@$(GO) tool cover -html=coverage.txt -o cover.html
	@open cover.html

## coverage-full: run go coverage test (with the long tests)
coverage-full: | $(GOBASE)
	@echo "  >  gocov ..."
	@$(CGO) test $(COVER_TEST_TARGETS) -v -race -coverprofile=coverage.txt -covermode=atomic -timeout=20m | tee coverage.log
	@$(GO) tool cover -html=coverage.txt -o cover.html
	@open cover.html

## codecov: run go test for codecov; (codecov.io)
codecov: | $(GOBASE)
	@echo "  >  codecov ..."
	@$(CGO) test $(COVER_TEST_TARGETS) -v -race -coverprofile=coverage.txt -covermode=atomic
	@bash <(curl -s https://codecov.io/bash) -t $(CODECOV_TOKEN)

## cyclo: run gocyclo tool
cyclo: | $(GOBASE) $(GOCYCLO)
	@echo "  >  gocyclo ..."
	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
	$(GOCYCLO) -top 20 .

## bench-std: benchmark test
bench-std:
	@echo "  >  benchmark testing ..."
	@$(CGO) test -bench="." -run=^$ -benchtime=10s $(COVER_TEST_TARGETS)
	# go test -bench "." -run=none -test.benchtime 10s
	# todo: go install golang.org/x/perf/cmd/benchstat


## bench: benchmark test
bench:
	@echo "  >  benchmark testing (manually) ..."
	@$(eval CGO_ENABLED = 1)
	@$(GO) test ./fast -v -race -run 'TestQueuePutGetLong' -timeout=20m


## linux-test: call ci/linux_test/Makefile
linux-test:
	@echo "  >  linux-test ..."
	@-touch $(STDERR)
	@-rm $(STDERR)
	@echo $(MAKE) -f ./ci/linux_test/Makefile test 2> $(STDERR)
	@$(MAKE) -f ./ci/linux_test/Makefile test 2> $(STDERR)
	@echo "  >  linux-test ..."
	$(MAKE) -f ./ci/linux_test/Makefile all  2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/' 1>&2





.PHONY: directories

directories: $(GOBIN)

MKDIR_P = mkdir -p
$(GOBIN):
	$(MKDIR_P) $(GOBIN)





.PHONY: printvars info help all
printvars:
	$(foreach V, $(sort $(filter-out .VARIABLES,$(.VARIABLES))), $(info $(v) = $($(v))) )
	# Simple:
	#   (foreach v, $(filter-out .VARIABLES,$(.VARIABLES)), $(info $(v) = $($(v))) )

print-%:
	@echo $* = $($*)

info:
	@echo "     GO_VERSION: $(GOVERSION)"
	@echo "        GOPROXY: $(GOPROXY)"
	@echo "         GOROOT: $(shell go env GOROOT) | GOPATH: $(shell go env GOPATH)"
	@echo "    GO111MODULE: $(GO111MODULE)"
	@echo
	@echo "         GOBASE: $(GOBASE)"
	@echo "          GOBIN: $(GOBIN)"
	@echo "    PROJECTNAME: $(PROJECTNAME)"
	@echo "        APPNAME: $(APPNAME)"
	@echo "        VERSION: $(VERSION)"
	@echo "      BUILDTIME: $(TIMESTAMP)"
	@echo "    GIT_VERSION: $(GIT_VERSION)"
	@echo "   GIT_REVISION: $(GIT_REVISION)"
	@echo "        GIT_HASH: $(GIT_HASH)"
	@echo "    GIT_SUMMARY: $(GIT_SUMMARY)"
	@echo "       GIT_DESC: $(GIT_DESC)"
	@echo
	@echo "             OS: $(OS)"
	@echo
	@echo " MAIN_BUILD_PKG: $(MAIN_BUILD_PKG)"
	@echo "      MAIN_APPS: $(MAIN_APPS)"
	@echo "       SUB_APPS: $(SUB_APPS)"
	@echo "MAIN_ENTRY_FILE: $(MAIN_ENTRY_FILE)"
	@echo
	#@echo "export GO111MODULE=on"
	@echo "export GOPROXY=$(shell go env GOPROXY)"
	#@echo "export GOPATH=$(shell go env GOPATH)"
	@echo

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
