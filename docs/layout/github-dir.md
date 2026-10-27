# Layout — `.github/`, `scripts/`, `tools/`, `templates/`, `examples/`

## `.github/`

```text
.github/
├── workflows/                  # 8 upstream workflows — DISABLED on the fork (NOIZU.md rule 6)
│   ├── ci.yaml                 # lint + unit/acceptance matrix
│   ├── codeql.yaml             # CodeQL scanning
│   ├── release.yaml            # goreleaser + SBOM + artifact attestations
│   ├── dotcom-acceptance-tests.yaml / ghes-acceptance-tests.yaml
│   └── labeler.yaml, stale.yaml, immediate-response.yaml
├── instructions/               # 5 path-scoped contributor instructions (go, docs, examples, schema-and-state, tests)
├── ISSUE_TEMPLATE/             # bug / documentation / feature / maintenance
├── CODEOWNERS, labeler.yml, dependabot.yml, release.yml
├── copilot-instructions.md     # repo-wide contributor guidance
└── pull_request_template.md
```

Workflows exist in the tree because the fork mirrors upstream; they are **not enabled** on `noizu-forks`. Nothing in this repo depends on CI running.

## `scripts/noizu/`

| File | Purpose |
|---|---|
| `build-dev-override.sh` | Builds the fork to `~/.local/share/opentofu/github-fork/terraform-provider-github` and writes `~/.config/opentofu/github-fork.tfrc` with `dev_overrides` for `registry.opentofu.org/integrations/github` (plus the terraform.io key). Refuses empty output dir and `/`. Prints commit + binary sha256. |

Use only with `TF_CLI_CONFIG_FILE` set — never installed globally.

## `tools/`

Separate Go module (`tools/go.mod`, `tools/go.sum`, `tools/tools.go`) that pins `tfplugindocs` behind a `//go:build generate` tag. The `go:generate` directives run `terraform fmt -recursive ../examples/` and `tfplugindocs generate --provider-dir ..`.

## `templates/`

Sources for generated docs: `index.md.tmpl`, `resources.md.tmpl`, `data-sources.md.tmpl`, plus per-type override directories (`templates/resources/`, `templates/data-sources/`). `make generatedocs` / `make checkdocs` regenerate `docs/`.

## `examples/`

21 top-level example directories: `provider/`, `resources/` (87), `data-sources/` (74), plus scenario examples (`app_authentication`, `app_token`, `emu`, `enterprise_settings`, `hosted_runner`, `release`, …) and `dev.tfrc`. Examples feed both doc generation and acceptance-test configs.

## `docs/`

Generated provider reference (`resources/` 88 files, `data-sources/` 75 files, `index.md`) **plus** hand-authored docs (`PROJ-LAYOUT.md`, `PROJ-ARCH.md`, `PROJ-SCHEMA.md`, `THREAT-MODEL.md`, their `layout/`, `arch/`, `schema/`, `threats/` subdirectories). Regeneration touches only the generated set.
