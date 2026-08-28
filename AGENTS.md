# AGENTS.md

Small Go web app (module `github.com/jwiklund/ah`) that records per-account yearly balances. Go>=1.27, no framework; plain `net/http` + `julienschmidt/httprouter`.

## Commands

- `make dev` — run dev server with hot-reload via `reflex`; serves on port 8080 against `example.txt`. Requires reflex first: `make install-reflex`.
- `make test` / `make test-w` — full suite (`go test ./...`); `-w` re-runs tests on file change.
- `make install` — `go install ./cmd/account-history`.
- `make update` — `go get -u ./... && go mod tidy`.

Gotchas:
- Makefile uses `.RECIPEPREFIX = >` (recipes indented with `>`), so it requires GNU Make 4.0+.
- `dev` watches `.go` files only; template changes are picked up because dev passes `-app.assets=view/assets` (lazy FS renderer). See View below.
- `make update` can bump `golang.org/x/exp` far ahead; check go.mod.

## Layout

- `cmd/account-history/main.go` — entrypoint, config, routes. Route list here is the source of truth for URLs.
- `control/` — HTTP handlers. One file per view (`index_control.go`, `edit_control.go`, etc.) plus shared helpers in `control.go`.
- `history/` — data model, YAML load/save, summary computation. All methods take the mutex internally.
- `csv/` — CSV/TSV import: separator/column type detection, plugin execution, row update application.
- `view/` — HTML templates + `//go:embed` assets.

## Config & data

- Config via `cristalhq/aconfig`: flags `-app.*`, env `APP_*`, and `$XDG_CONFIG_HOME/account-history/config.yaml`. Defaults: `accounts.txt`, port 8080.
- Plugin config `plugins.yaml` (same config dir), `import:` key, `KnownFields(true)`.
- Accounts storage is a **multi-document YAML stream** — accounts separated by `---` lines; each has `name`, `history` (list of `date`/`amount`/`change`), optional `tags`. Save copies the file to `<file>.bak` first.
- CLI extras: `-init-history` creates a file if missing, `-add-year=YYYY` adds a year.

## View / templates

- `view/assets/assets.go` embeds `*.html *.ico` only. Adding e.g. `.css`/`.js` files there will NOT be embedded when running without `-app.assets`.
- Templates render full pages and htmx partials from the same file: the `{{block "page.body.html" .}}` block is what htmx requests (`Hx-Request` header) render, via `control/templateNameWithPart`. New page => keep the `page.body.html` block name in sync.
- Template funcs: `human` (thousand separators), `json`. `view.New` renders `index.html` once at startup to catch template errors.

## Conventions

- `history.NameToSlug(name)` lowercases and collapses non-alphanumerics to `-`; used for route params and dedup.
- Tag `Oneoff` (const in `history/types.go`) changes summary/sort semantics; don't rename strings.
- Uses `golang.org/x/exp/slices` (not stdlib `slices`) and `stretchr/testify`.
- Tests: `history/*_test.go`, `csv/*_test.go`; csv fixtures live in `csv/assets/*.yaml`.