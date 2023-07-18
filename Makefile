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

.PHONY: all
all: dev

mh: $(find . -name "*.html" -o -name "*.go")
> go build -o mh cmd/mh/main.go

.PHONY: dev
dev: $(shell bin/has reflex)
> reflex -d none -s -R vendor. -r \.go$$ -- go run cmd/mh/main.go -app.assets=view/assets -app.accounts=example.txt

.PHONY: test
test:
> go test ./...

.PHONY: test-w
test-w: $(shell bin/has reflex)
> reflex -d none -s -R vendor. -r \.go$$ -- go test ./...

install-reflex:
> go install github.com/cespare/reflex@latest && touch install-reflex
