SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

DEV_ACCOUNTS := -app.accounts=example.txt

.PHONY: all
all: dev

install: $(shell find . -name "*.html" -o -name "*.go")
> go install ./cmd/account-history

.PHONY: dev
dev: $(shell bin/has reflex)
> reflex -d none -s -R vendor. -r \.go$$ -- go run cmd/account-history/main.go -app.assets=view/assets $(DEV_ACCOUNTS)

.PHONY: test
test:
> go test ./...

.PHONY: test-w
test-w: $(shell bin/has reflex)
> reflex -d none -s -R vendor. -r '\.go$$|\.yaml$$' -- go test ./...

.PHONY: update
update:
> go get -u ./...
> go mod tidy

install-reflex:
> go install github.com/cespare/reflex@latest
