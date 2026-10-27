# NOIZU.md — noizu-forks/terraform-provider-github

<!-- 〚📁:b2de3d9e-bdf7-5743-9e73-6c02a33041f7〛 Noizu Fork Contract :: Fork contract: how the noizu extension layer stays conflict-free against upstream. -->

Noizu's fork of [integrations/terraform-provider-github](https://github.com/integrations/terraform-provider-github),
carrying an extension layer for resources the upstream provider does not
(plan to) ship. This file is the contract for how the fork stays
conflict-free against upstream.

## Purpose

Upstream covers standard GitHub resources well, but some GitHub App /
org-admin surfaces are missing (app webhook configuration via JWT,
code security configurations, copilot seat management, …). Rather than
patching upstream files, all Noizu additions live behind a tiny
registration hook so upstream pulls stay merge-clean.

## Conflict-avoidance rules (the contract)

1. **Upstream files get exactly ONE touch, ever.** At v6.13.0 that touch
   is in `github/provider.go`: the return expression of `NewProvider` is
   wrapped — `return registerNoizuExtensions(&schema.Provider{...})` —
   because the provider literal is returned directly from the closure
   (there is no named `p` to mutate before a bare return). Two physical
   lines (the `return &schema.Provider{` opener and its closing brace),
   one logical edit site. Everything else Noizu does must be NEW files.
2. **Noizu code lives in `github/noizu_*.go`, package `github`.** Same
   package on purpose — avoids import cycles with the provider's
   `Config`/`Owner` helpers and lets extensions reuse `ghclient`.
3. **Naming convention.** Resource names `github_noizu_<thing>`;
   files `github/noizu_resource_<thing>.go`,
   `github/noizu_data_source_<thing>.go`, tests alongside
   (`*_test.go`). Register in `noizuResources` / `noizuDataSources` maps
   in `github/noizu_provider.go`.
4. **Collisions panic.** `registerNoizuExtensions` panics if a noizu name
   already exists upstream, so drift fails loudly at provider
   construction, never at plan time. Guarded by unit tests.
5. **Changelog.** Upstream `CHANGELOG.md` is never edited. At v6.13.0
   upstream carries no root `CHANGELOG.md` anyway (release notes live in
   `.changelog/`); Noizu changes are logged in `NOIZU-CHANGELOG.md`.
6. **CI.** GitHub Actions stay disabled on the fork; never enable them.
   Acceptance tests are never run by agents (they hit live GitHub) —
   unit tests only (`TF_ACC` unset).

## Branch model

- `main` — pristine mirror of upstream `main`. Never commit to it.
- `develop` — the Noizu line, cut from an upstream **release tag**
  (currently `v6.13.0`). All work lands here via PR (squash).
- **Upgrading:** when upstream publishes a new release tag, merge the tag
  into `develop` (or re-cut if history is trivial). The single-touch rule
  keeps that merge conflict-free or near-zero.
- **Generic bug fixes** to upstream resources go upstream as PRs to
  integrations/terraform-provider-github; while carried here pre-merge,
  cherry-pick them onto `develop` with commit subjects prefixed
  `[upstream-pr #N]` so they are identifiable and droppable once
  upstreamed.

## Planned gap resources (design TBD, not yet implemented)

Terraform lane (via the extension layer):
- GitHub App webhook configuration (App API via JWT)
- Code security configurations (org/repo attach)
- Copilot seat management (org)

Browser-lane only (UI automation, not provider resources): app
permissions/events/URLs, key rotation, copilot policies.

## Consuming the fork (dev_overrides)

```bash
./scripts/noizu/build-dev-override.sh          # build + write ~/.config/opentofu/github-fork.tfrc
export TF_CLI_CONFIG_FILE="$HOME/.config/opentofu/github-fork.tfrc"
tofu init && tofu plan
```

The script builds to `~/.local/share/opentofu/github-fork/terraform-provider-github`
and writes a dedicated CLI config whose `dev_overrides` primary key is
`registry.opentofu.org/integrations/github` (OpenTofu address; the
terraform.io address is included for `terraform`-CLI compatibility).
Use it only with `TF_CLI_CONFIG_FILE` set — never as a global config.
