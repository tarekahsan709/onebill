# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

**OneBill** — a self-hosted cloud + AI spend dashboard, and the skeleton for it. Originally a fork of [Miniflux 2](https://github.com/miniflux/v2) (Go feed reader, Apache 2.0). The full product plan lives in `cloud-cost-dashboard-spec.md`, which is gitignored and local-only — read it if it is present, otherwise the summary below is the plan.

**Step 1 of the plan is done:** everything feed-reader specific was deleted, and the app was renamed from Miniflux to OneBill. The git remote was removed; this repo no longer tracks upstream, so there is nothing to merge from and no reason to preserve upstream compatibility. What remains is roughly 22k lines of Go: users, sessions, auth, settings, a REST API, migrations, config, CLI, and Docker packaging.

The shape being reused: poll external sources on a schedule, normalize into Postgres, serve a small server-rendered UI plus a REST API. Swap "RSS feed" for "cloud provider" and the architecture holds.

Module path is `onebill.app`, binary is `onebill`, default database is `onebill`.

**The domain model still uses feed-reader nouns.** `feeds` is the provider table and `entries` is the cost-record table. Renaming those is step 2 (below) — until then, read "feed" as "provider" and "entry" as "cost record".

### Licensing

OneBill is Apache 2.0, the same as its upstream. Section 4 of that license requires keeping the origin's attribution in derivative works, so:

- Every Go file carries **both** copyright lines — Miniflux's first, OneBill's second. `.golangci.yml`'s `goheader` template enforces the exact three-line block; new files must match it.
- [NOTICE](NOTICE) credits Miniflux and must ship with any distribution.
- The About page credits string in `en_US.json` names both.

Do not strip the Miniflux copyright lines. Everything else that said "Miniflux" is gone.

## Commands

```bash
make onebill                  # build binary for the current platform
make run                      # run locally: debug logs, migrations, admin/test123
make test                     # go test -cover -race -count=1 ./...
make lint                     # go vet + gofmt check + golangci-lint
make integration-test         # spin up postgres + server, run ./internal/api against it
make clean-integration-test   # kill the server and drop onebill_test
make docker-image             # packaging/docker/alpine/Dockerfile

go test -run TestValidateEntryStatus ./internal/validator/   # single test
go test -race ./internal/storage/                            # single package
```

Needs Go >= 1.26 and PostgreSQL. Local database for development:

```bash
docker run --rm --name onebill-db -p 5432:5432 \
  -e POSTGRES_DB=onebill -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres postgres
```

`make run` uses the default `DATABASE_URL` (`user=postgres password=postgres dbname=onebill sslmode=disable`) and serves on `127.0.0.1:8080`.

JavaScript in [internal/ui/static/js/](internal/ui/static/js/) is linted separately in CI with `jshint` and `eslint` — no build step, no bundler.

## Architecture

**Single binary, single process.** [main.go](main.go) calls `cli.Parse()`, which parses flags, loads config from the environment, opens Postgres, and either runs a one-shot task (`-migrate`, `-create-admin`, `-run-cleanup-tasks`, …) or starts the daemon.

**The daemon** ([internal/cli/daemon.go](internal/cli/daemon.go)) starts the HTTP server plus a cleanup ticker, then blocks on signals. **There is no worker pool or sync scheduler right now** — the feed-fetching machinery went with `internal/reader`, and the provider sync (step 4 below) will bring back the same shape: a ticker that selects providers whose `next_check_at` has expired, and a pool of goroutines that call each provider's billing API and write rows.

The ingestion seam that does exist today is `POST /v1/feeds/{feedID}/entries/import` → `handler.importFeedEntryHandler` → `store.InsertEntryForFeed`. That is the closest thing to "write a cost record" and worth reading before writing the sync.

**Routing** is plain `net/http.ServeMux` with Go method+path patterns (`"GET /feed/{feedID}/edit"`). Two levels: [internal/http/server/routes.go](internal/http/server/routes.go) mounts the API (`/v1/`), metrics, and the UI catch-all, then wraps them in shared middleware; the outer router strips `BASE_PATH` so every inner handler sees a clean path. Health probes are registered at the true root so they ignore the base path.

**UI** ([internal/ui/](internal/ui/)) — server-rendered HTML, one file per page, each a method on the shared `handler` struct holding `basePath`, `store`, and `tpl`. The pattern is always: load user → call storage → `view.New(h.tpl, r).Set(...).Render("template_name")` → `response.HTML(...)`. Auth, sessions, and CSRF are middleware in the same package; [internal/ui/routes.go](internal/ui/routes.go) declares which paths are public.

**Templates** ([internal/template/](internal/template/)) — `html/template`, embedded with `go:embed`. Every view under `templates/views/` **must** be listed in `Engine.ParseTemplates()` with its layout dependencies, or the app panics at boot. Templates are parsed once, then cloned per request so the per-language `t`/`plural` funcs don't race. Template field errors only surface at render time, so click through a page after changing a model.

**Static assets** ([internal/ui/static/](internal/ui/static/)) — CSS, JS, and icons are embedded, minified at startup, and served under checksum paths (`/js/{checksum}/{filename}`) for cache busting. `GenerateJavascriptBundles` lists the files in each bundle by hand; deleting a `.js` file means editing that list too.

**Storage** ([internal/storage/](internal/storage/)) — hand-written SQL on `database/sql`, no ORM. One file per entity. Complex reads use small builder types (`entry_query_builder.go`, `feed_query_builder.go`) that accumulate `conditions` and positional `$N` args. Every query is scoped by `user_id`; keep it that way, and add `organization_id` alongside it in step 2.

**Migrations** ([internal/database/migrations.go](internal/database/migrations.go)) — an ordered array of `func(tx *sql.Tx) error`, currently 135 entries. `schemaVersion` is just `len(migrations)`. **Append only. Never edit or reorder an existing migration** — the version number is the index. Migration 135 is the one that dropped every table and column the deleted features owned.

**Config** ([internal/config/](internal/config/)) — environment variables parsed once at boot into the global `config.Opts`, with `Validate()` rejecting bad combinations. Read config through `config.Opts`; do not read `os.Getenv` elsewhere. Options for deleted subsystems (fetcher, media proxy, YouTube, polling scheduler) were removed along with their man page entries.

**Errors shown to users** go through `locale.LocalizedError` — a translation key plus args, rendered at the point of display via `Translate(language)`. Internal errors stay plain `error`s.

**i18n** ([internal/locale/](internal/locale/)) — `en_US` only now. `make add-string KEY=my.key VAL="English text"` still works and writes to every catalog in `translations/`. `TestAllKeysHaveValue` fails on an empty value.

## Conventions

- Every Go file starts with the SPDX header — the `goheader` linter enforces it:
  ```go
  // SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
  // SPDX-FileCopyrightText: Copyright The OneBill Authors. All rights reserved.
  // SPDX-License-Identifier: Apache-2.0
  ```
- Package clauses carry an import comment: `package ui // import "onebill.app/internal/ui"`.
- Commit messages must match conventional commits; CI checks every commit in a PR (`.github/workflows/scripts/commit-checker.py`). Allowed types: `build chore ci docs feat fix perf refactor revert security style test`, optional lowercase scope.
- `errcheck` is off; `gocritic`, `staticcheck`, `perfsprint`, `sqlclosecheck`, `loggercheck`, `misspell` are on.
- Logging is `log/slog` with structured attributes (`slog.Int64("feed_id", …)`), never `fmt` into a message.
- Dependencies are kept deliberately few. Adding one needs a real reason.
- Vanilla JS only, no framework, no React.

## The plan

One page showing total cloud + AI spend this month, last month, and a breakdown by provider. Nothing else — no forecasting, no anomaly detection, no rightsizing, no chargeback. The wedge is self-hosted (credentials never leave the customer's infrastructure), European hosters that no commercial tool supports (Hetzner, OVH, Scaleway, IONOS), and simplicity.

Step 1 (strip the feed reader, rename to OneBill) is done. Remaining, in order:

2. **Rename the domain model**: `feeds` → `providers` and `entries` → `cost_records` via new migrations, plus the matching Go types, handlers, routes, and templates. Add `organization_id` from the start and apply it to every query alongside `user_id`. Add the provider columns the spec calls for: `kind`, `credentials`, `currency`.
3. **`internal/providers/`** — one file per provider behind a `Provider` interface (`Kind`, `Validate`, `Fetch`). Adding a provider should mean adding one file.
4. **Bring back the worker pool and scheduler**, retargeted: tick → select due providers → fetch billing data → write cost records. Same shape as the deleted `internal/worker`; it's in git history if you want the original.
5. **Credential encryption** at rest, key from the environment, never logged. The genuinely new hard part — [internal/crypto/](internal/crypto/) has only hashing and random helpers today.
6. **Currency**: store the original amount and the converted amount, plus the FX rate and its date. Totals get disputed otherwise.
7. **Monthly aggregation** (`SUM(amount_cents_base) GROUP BY provider, month`) replacing the entry list.

Three pages total: dashboard, providers, settings. v1 covers Hetzner, AWS, Anthropic with a daily sync. No forecasting, alerts, or billing code.

Hetzner is the checkpoint — if its number matches the real invoice, the rest is repetition.

## Known stale after the deletion pass

- `github.com/onebill/onebill` is a placeholder in the Docker labels, Debian control file, and RPM spec, and the CI workflows gate on `github.repository_owner == 'onebill'`. Set these once the real repository exists, or those release jobs will never fire.
- `make integration-test` starts a real server but `./internal/api` has only unit tests now; the end-to-end API tests went with the deleted Go client.

Everything else was swept: `golang.org/x/tools/cmd/deadcode -test ./...` reports nothing, CSS rules and custom properties with no markup behind them are gone, and [onebill.1](onebill.1) no longer documents removed flags or env vars. Re-run `deadcode` after any large deletion — grep-based audits give false positives on interface methods and intra-package calls.
