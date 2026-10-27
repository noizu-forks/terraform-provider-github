# Architecture — Noizu Extension Layer

The fork's whole reason for existing: ship resources upstream does not have, without diverging from upstream code.

## The contract (`NOIZU.md`)

1. **Upstream files get exactly ONE touch, ever.** In this tree that touch is `github/provider.go`: `NewProvider`'s closure returns `registerNoizuExtensions(&schema.Provider{...})`. Everything else Noizu adds must be a new file.
2. **Same package, `noizu_*` filenames.** `github/noizu_*.go`, package `github`, so extensions can use `Config`, `Owner` and `ghclient` without import cycles.
3. **Namespaced names.** Resources `github_noizu_<thing>`; files `noizu_resource_<thing>.go` / `noizu_data_source_<thing>.go`; tests alongside.
4. **Collisions panic.** Registration fails at provider construction, never at plan time.
5. **Changelog.** `NOIZU-CHANGELOG.md` only; upstream `CHANGELOG.md` is never edited.
6. **CI stays disabled.** Acceptance tests are never run by agents (they hit live GitHub); unit tests only.

## Mechanism

```go
var noizuResources  = map[string]*schema.Resource{}   // github_noizu_* managed resources
var noizuDataSources = map[string]*schema.Resource{}  // github_noizu_* data sources

func registerNoizuExtensions(p *schema.Provider) *schema.Provider {
    mergeNoizuMaps("resource", p.ResourcesMap, noizuResources)
    mergeNoizuMaps("data source", p.DataSourcesMap, noizuDataSources)
    return p
}
```

`mergeNoizuMaps` panics with an actionable message if a name already exists upstream. `github/noizu_provider_test.go` covers the merge and the collision panic, which is why the guard is trustworthy rather than aspirational.

## Current state (honest)

Both maps are **empty**. The 2026-10-14 scaffold registered nothing; "planned gap resources" in `NOIZU.md` (App webhook configuration via JWT, code security configurations, Copilot seat management) are design-only, and the browser-lane items are explicitly *not* provider resources. Until a `noizu_resource_*.go` lands, the fork behaves exactly like upstream v6.13.0.

## Consuming the fork

```bash
./scripts/noizu/build-dev-override.sh
export TF_CLI_CONFIG_FILE="$HOME/.config/opentofu/github-fork.tfrc"
tofu init && tofu plan
```

`dev_overrides` keys: `registry.opentofu.org/integrations/github` (canonical) and `registry.terraform.io/integrations/github` (CLI compatibility). Keep it scoped to `TF_CLI_CONFIG_FILE`; making it global would shadow the real provider for every project on the machine.

## Upgrading upstream

Merge the new upstream release tag into `develop`. The single-touch rule is what makes that merge trivial; a conflict in `github/provider.go` means something else touched an upstream file and should be re-homed into a `noizu_*.go` file.
