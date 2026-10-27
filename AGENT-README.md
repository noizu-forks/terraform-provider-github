# AGENT-README

> Orientation for AI coding agents working in this repo. Human docs: README.md (upstream) and NOIZU.md (fork contract — read it first).

## Type

terraform+go

## Purpose

Noizu fork of integrations/terraform-provider-github carrying an extension layer (`github_noizu_*` resources/data sources) for GitHub App / org-admin surfaces upstream does not ship (app webhook config via JWT, code security configurations, copilot seats, …).

## Stack

Go Terraform provider (upstream v6.13.0 base). golangci-lint gates (base + strict rule files merged into `.golangci.new.yml`).

## Dev mode

```bash
make build                              # lintcheck + CGO_ENABLED=0 go build ./...
scripts/noizu/build-dev-override.sh     # dev_overrides install tooling
```

## Ports / URLs

- none (library/CLI)

## Tests

```bash
make test          # unit tests only; T=<pattern> filters, COV=true adds coverage
```

## Gotchas

- Fork contract (NOIZU.md): upstream files get exactly ONE touch ever (`NewProvider`'s return wraps `registerNoizuExtensions`); all Noizu code lives in new files `github/noizu_*.go`, names `github_noizu_*`; collisions panic at provider construction.
- Acceptance tests hit live GitHub — never run by agents; unit tests only (`TF_ACC` unset).
- GitHub Actions stay disabled on this fork — never enable them.
- Branch model: `main` is a pristine upstream mirror (never commit to it); all work lands on `develop` via PR (squash).
- Upstream `CHANGELOG.md` is never edited; Noizu changes are logged in `NOIZU-CHANGELOG.md`.
