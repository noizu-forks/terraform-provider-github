# Threat — Supply Chain and Build Provenance

Covers T-001, T-004, T-006, T-008, T-009, T-012.

## Fork build is unsigned and shadowing (T-001, High, Partial)

`scripts/noizu/build-dev-override.sh` compiles the working tree to `~/.local/share/opentofu/github-fork/terraform-provider-github` and writes a CLI config whose `dev_overrides` map `registry.opentofu.org/integrations/github` **and** `registry.terraform.io/integrations/github` to that directory.

Consequences:

- Terraform will run whatever that path contains, with no signature, checksum or attestation check. Upstream's SLSA/cosign attestations (`VERIFY_ATTESTATIONS.md`) do **not** apply to fork builds.
- Any project that inherits `TF_CLI_CONFIG_FILE` (shell profile, CI variable, `.envrc`) silently uses the fork. A stale or tampered binary is indistinguishable from a good one at plan time.

Controls present: the config is a dedicated file used only via `TF_CLI_CONFIG_FILE` (never a global CLI config); the script refuses an empty output dir or `/` and prints commit + binary sha256; `NOIZU.md` documents the scoping rule. There is no integrity gate — the sha256 is informational.

## No CI on the fork (T-009, Medium, Open)

Workflows exist in `.github/workflows/` but are disabled by fork policy. Practically:

- No lint, unit test, CodeQL or dependency-advisory run on push. `go.sum` pinning plus manual review of upstream dependency bumps is the only control.
- Because `make build` depends on `lintcheck`, a developer who builds locally does get lint; nothing forces it.

## Dependency surface

Heavy third-party reliance: `go-github/v88`, `githubv4`, `terraform-plugin-sdk/v2`, `terraform-plugin-testing`, `github-conditional-http-transport` (+ bbolt), `go-github-ratelimit/v2`, `go-githubauth`, `go-jose/v4`. The conditional-request cache library writes the on-disk store, so a cache-library defect is also a local-file defect.

## Response cache (T-004, Medium, Partial)

`<cache_path>/terraform-provider-github/<ref>/cache.db` persists API response bodies (repo/org metadata, not credentials). Permissions are `0700` dir / `0600` file, which contains the blast radius on a single-user host and does not help on a shared runner where the same UID is used by other jobs.

## Destructive authority (T-006, High, Accepted)

The provider is authorised to delete repositories, teams, and org-level settings. A wrong `owner`, a `base_url` aimed at the wrong enterprise, or a stale state produces destruction with no provider-side guard, dry-run, or confirmation. Consumer-side controls: state locking, `prevent_destroy`, plan review. Nothing in this repo substitutes for them.

## Rate-limit amplification (T-008, Low, Mitigated)

`parallel_requests` defaults off; read/write delays default 0/1000 ms; `max_retries` defaults 3 against 500/502/503/504 only. A pathological retry loop is bounded. In the new client, delays are library-managed and `parallel_requests` is ignored.

## Extension-layer collision guard (T-012, Low, Mitigated)

`mergeNoizuMaps` panics if a Noizu name already exists upstream, or vice versa, at provider construction. This prevents a scenario where a trusted resource name silently resolves to different code. Guarded by `github/noizu_provider_test.go`.

## Residual

Fork consumption is a **trust-the-working-tree** model. If that is not acceptable for a given environment, the alternative is to build in an isolated, reviewed checkout and compare the printed sha256 against a recorded value out of band — which is a process control, not a technical one.
