# Architecture — Decision Notes

Condensed from `DECISIONS.md` (which stays authoritative), plus the fork-level decisions that only exist in `NOIZU.md`.

## Upstream decisions still in force

| Decision | Why it matters when changing code |
|---|---|
| One resource = one API entity | Do not merge entities into a resource; state and destroy semantics depend on the 1:1 mapping |
| Minimize API calls | Single read populating all computed fields; never call Read at the end of Create/Update, never Read after Delete |
| API-only, no local git | Resources must not shell out to `git`; all git operations go through the API |
| Explicit authentication | Auth comes only from documented config/env/`gh` sources — no ambient discovery |
| Transport rework | Conditional requests + rate limiting are library-backed in the new client; don't hand-roll a third stack |
| `StateUpgraders`, not `MigrateState` | New schema changes ship a `resource_github_<entity>_migration.go` + `SchemaVersion` bump |
| `terraform-plugin-testing` | New tests use `ConfigStateChecks`/`ConfigPlanChecks` |
| `tflog`, not `log` | Structured logging only in new code |
| Generated docs | Documentation is generated from `templates/` + `examples/` via tfplugindocs |
| SDK v2 | Context-aware CRUD, `ValidateDiagFunc`, `diag.Diagnostics`, `any` |

## Fork-level decisions (this tree only)

| Decision | Rationale | Cost |
|---|---|---|
| Single upstream touch in `provider.go` | Upstream release merges stay near-conflict-free | The return-expression wrapper is unusual; do not "clean it up" |
| Same package (`github`) for `noizu_*` files | Access `Config`/`Owner`/`ghclient` without import cycles | No package boundary; enforce via naming + review |
| `github_noizu_` name prefix + panic on collision | Drift fails at construction, not at plan time | Upstream must never ship a `github_noizu_*` name (it won't) |
| `NOIZU-CHANGELOG.md` instead of upstream changelog | Zero conflict surface in a file upstream rewrites per release | Fork history lives in a second file |
| CI disabled on the fork | Actions on a fork against upstream credentials is a leak risk and unwanted noise | No automated gate; local `make test` / `make build` is the contract |
| `dev_overrides` consumption | No registry publish, no image, no signed artifact | Binary provenance is local; see `THREAT-MODEL.md` |
| Modules keep upstream's module path | Avoids rewriting every import and the registry address | `go.mod` looks like upstream; fork identity is in the remote |

## Open design questions (recorded, not decided)

- Whether the three planned gap resources (App webhook config, code security configs, Copilot seats) belong in the Terraform lane at all — for the App API surfaces the JWT flow and the org-level scoping need a design pass first.
- Whether `legacy_client` should default to `false` in the fork or stay aligned with upstream (staying aligned is the current, deliberate choice).
