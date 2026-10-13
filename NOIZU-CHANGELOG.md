# Noizu fork changelog

Noizu changes only. Upstream `CHANGELOG.md` is never edited (conflict
avoidance — see `NOIZU.md`); at v6.13.0 upstream carries no root
CHANGELOG.md anyway.

## 2026-10-14

- Scaffolded the Noizu extension layer: `github/noizu_provider.go`
  (`noizuResources` / `noizuDataSources` maps, collision-panicking
  `registerNoizuExtensions`), `github/noizu_provider_test.go` unit
  tests, `scripts/noizu/build-dev-override.sh` dev_overrides tooling,
  and this file. Single upstream touch: `NewProvider`'s return
  expression in `github/provider.go` now wraps
  `registerNoizuExtensions(...)`. No resources registered yet. (PR #1)
