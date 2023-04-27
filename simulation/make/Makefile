###
# Find OS and Go environment
# GO contains the Go binary
# FS contains the OS file separator
###
ifeq ($(OS),Windows_NT)
  GO := $(shell where go.exe 2> NUL)
  FS := "\\"
else
  GO := $(shell command -v go 2> /dev/null)
  FS := "/"
endif

ifeq ($(GO),)
  $(error could not find go. Is it in PATH? $(GO))
endif

GOPATH ?= $(shell $(GO) env GOPATH)
GITHUBDIR := $(GOPATH)$(FS)src$(FS)github.com

###
# tools
###

TOOLS_DESTDIR  ?= $(GOPATH)/bin

RUNSIM			  = $(TOOLS_DESTDIR)/runsim

all: tools

tools: runsim


# Install the runsim binary with a temporary workaround of entering an outside
# directory as the "go get" command ignores the -mod option and will pollute the
# go.{mod, sum} files.
# 
# ref: https://github.com/golang/go/issues/30515
runsim: $(RUNSIM)
$(RUNSIM):
	@echo "Installing runsim..."
	@(cd /tmp && go install github.com/cosmos/tools/cmd/runsim@v1.0.0)

tools-clean:
	rm -f $(RUNSIM)
	rm -f tools-stamp

.PHONY: all tools tools-clean
