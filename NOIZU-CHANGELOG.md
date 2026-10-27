# Noizu fork changelog

<!-- 〚📁:f3260096-6a15-5bc3-994f-700eb8e8bcf1〛 Noizu Fork Changelog :: Dated log of Noizu-only changes to the fork. -->

Noizu changes only. Upstream `CHANGELOG.md` is never edited (conflict
avoidance — see `NOIZU.md`); at v6.13.0 upstream carries no root
CHANGELOG.md anyway.

## 2026-10-27

- Added hand-authored project docs: `docs/PROJ-LAYOUT.md` (+`layout/`),
  `docs/PROJ-ARCH.md` (+`arch/`), `docs/PROJ-SCHEMA.md`,
  `docs/THREAT-MODEL.md` (+`threats/`), each with a `.summary.md`. Describes
  this fork tree only (provider schema, auth chain, transport, extension
  layer, on-disk stores, threat register). Two open findings recorded:
  `token` is not `Sensitive` in the provider schema, and the fork has no
  CI/dependency-advisory signal.

## 2026-10-14

- Scaffolded the Noizu extension layer: `github/noizu_provider.go`
  (`noizuResources` / `noizuDataSources` maps, collision-panicking
  `registerNoizuExtensions`), `github/noizu_provider_test.go` unit
  tests, `scripts/noizu/build-dev-override.sh` dev_overrides tooling,
  and this file. Single upstream touch: `NewProvider`'s return
  expression in `github/provider.go` now wraps
  `registerNoizuExtensions(...)`. No resources registered yet. (PR #1)
