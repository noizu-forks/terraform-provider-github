# Architecture — Testing, Build and Release

## Test tiers

| Tier | Command | Needs credentials | Notes |
|---|---|---|---|
| Unit | `make test` (`T=` filters) | no | runs everywhere, including agents |
| Acceptance | `make testacc T=TestAcc...` | yes (`TF_ACC`) | hits live GitHub; **never run by agents** |
| Sweep | `make sweep` (`SWEEP=repositories,teams`) | yes | deletes leaked `tf-acc-test-*` resources |

Acceptance gating lives in `github/acc_test.go`: `providerFactories`, `TF_ACC`, and mode skips (`skipUnauthenticated`, `skipUnlessHasOrgs`, `skipUnlessHasPaidOrgs`, `skipUnlessEnterprise`, `skipUnlessHasAppInstallations`, `skipUnlessEMUEnterprise`, `skipIfEMUEnterprise`). Modes are anonymous / individual / organization / enterprise, driven by `GITHUB_TOKEN` plus `GITHUB_OWNER` or `GITHUB_ORGANIZATION`.

## Assertion style

Modern tests use `ConfigStateChecks` (`statecheck`, `knownvalue`, `tfjsonpath`, `compare`) and `ConfigPlanChecks` (`plancheck`) from `terraform-plugin-testing`. Legacy `Check:` + `resource.ComposeTestCheckFunc` still exists in older tests — do not copy it. One `TestAcc...` function per resource with `t.Run` subtests (`ARCHITECTURE.md`).

## Build

- `make build` depends on `lintcheck`; `make lintcheck-strict` merges `.golangci.yml` + `.golangci.strict.yml` into a generated `.golangci.new.yml` (gitignored) and runs the stricter set — expected to fail on legacy code, used to stop new debt.
- Module path is upstream's (`github.com/integrations/terraform-provider-github/v6`) and Go 1.26 is required. Do not rename the module: resources, tests and the registry address all assume it.
- `main.go` reads `version`/`commit` from `-ldflags`, set by GoReleaser (`.goreleaser.yml`).

## Docs pipeline

`tools/tools.go` pins `tfplugindocs`; `make generatedocs` renders `docs/resources/`, `docs/data-sources/`, `docs/index.md` from `templates/` + `examples/`. `make checkdocs` is the drift gate. Validation/format/lint of docs: `make validatedocs`, `make fmtdocs`, `make lintdocs`, `make mdfmt`, `make mdlint`.

Generated files under `docs/resources/` and `docs/data-sources/` are outputs — edit `templates/` or `examples/` instead.

## Release

Upstream releases via `.github/workflows/release.yaml` (goreleaser + SBOM + artifact attestations; verification steps in `VERIFY_ATTESTATIONS.md`). **On this fork those workflows are disabled and no release pipeline runs.** The fork is consumed from a locally built binary through `dev_overrides`; there is no signed-artifact path for fork builds, so provenance guarantees apply to upstream artifacts only.
