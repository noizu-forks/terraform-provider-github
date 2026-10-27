# Project Architecture (summary)

Single Go plugin binary (Terraform/OpenTofu provider) — no server, database, or deploy manifests. This tree is Noizu's fork of `integrations/terraform-provider-github`, base tag `v6.13.0`; `main` mirrors upstream, `develop` carries the extension layer.

## Components

| Component | Purpose | Location |
|---|---|---|
| Plugin entry | `plugin.Serve`, addr `registry.terraform.io/integrations/github` | `main.go` |
| Provider schema | config attributes, 88 resources, 75 data sources, `Configure` auth chain | `github/provider.go` |
| Resources / data sources | context-aware CRUD, importers, state upgraders | `github/resource_*.go`, `github/data_source_*.go` |
| Noizu extension layer | `noizuResources`/`noizuDataSources` maps merged at construction, panic on collision (currently empty) | `github/noizu_provider.go` |
| Legacy client | `Config`/`Owner`, etag + rate-limit + retry transports | `github/config.go`, `github/transport.go` |
| New client | `Source` abstraction (token/app/anonymous), library-backed conditional requests + rate limit, bbolt cache | `internal/ghclient/` |
| Utilities | ID build/parse, validators, team/repo resolution, diff funcs | `github/util*.go` |

## Auth chain

`app_auth`/`GITHUB_APP_*` → `token`/`GITHUB_TOKEN` → `gh auth token` (legacy path only) → anonymous. Authenticated new-client use requires explicit `owner`. `legacy_client` (default `true`) selects the client stack; delay/retry args apply to the legacy stack.

## Flow

Provider construction registers upstream + Noizu maps → `Configure` resolves credentials and builds REST/GraphQL clients into `Owner` → resources issue targeted API calls → 404 in Read clears the ID (drop from state) → import sets minimum state, Terraform re-calls Read.

## Testing / build

Unit tests credential-free and agent-safe; acceptance tests hit live GitHub and are agents-forbidden. `make build` implies `lintcheck`; generated docs come from `templates/` + `examples/` via tfplugindocs. CI workflows exist but are disabled on the fork; consumption is a local binary through `dev_overrides`.

## Key decisions

One resource = one API entity; API-only (no local git); explicit auth; SDK v2 with `StateUpgraders` and `terraform-plugin-testing`; `tflog`; exactly one upstream touch (`NewProvider` return expression) with all Noizu code in new `noizu_*.go` files.

Details: [PROJ-ARCH.md](PROJ-ARCH.md), [arch/auth-and-clients.md](arch/auth-and-clients.md), [arch/transport.md](arch/transport.md), [arch/noizu-extension-layer.md](arch/noizu-extension-layer.md), [arch/testing-and-build.md](arch/testing-and-build.md), [arch/decisions.md](arch/decisions.md).
