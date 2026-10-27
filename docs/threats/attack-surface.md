# Threat — Attack Surface

## Ingress (untrusted-or-foreign input)

| Entry point | Who controls it | Notes |
|---|---|---|
| Provider block + `GITHUB_*` env vars | operator / CI config | `base_url`, `owner`, tokens, delays |
| Terraform config + state | operator / state backend | state is an input on every run and is plaintext |
| `gh` CLI subprocess | local machine | token on stdout; binary resolved via `GH_PATH`/`PATH` |
| GitHub API responses | GitHub (remote, semi-trusted) | drives Read/Computed state; a malicious/compromised response can poison state |
| Upstream source merged into `develop` | upstream project | code review is the only gate (no CI) |
| `.meta/pointers.yaml` | repo content | parsed by tooling, not at runtime |

## Egress (what leaves the process)

| Destination | Payload | Control |
|---|---|---|
| GitHub API | token in `Authorization`, plus resource payloads | TLS; base URL from config |
| GitHub API | app JWT → installation token | short-lived token exchange |
| Local disk | cache bodies, PEM reads | fs permissions |
| Local process table / logs | env-derived credentials | tflog debug messages only — verify no credential is logged |

## Stores

| Store | Contents | Protection |
|---|---|---|
| Terraform state | all resource attributes incl. secrets | consumer-chosen backend; **no encryption by provider** |
| `cache.db` (bbolt) | API response bodies | `0700` dir / `0600` file |
| `pem_file` | app private key | operator-managed file |
| CI env / process env | token, `GITHUB_APP_PEM_FILE` | operator-managed |

## Boundary crossings to reason about

1. **Operator → provider**: misconfiguration (wrong `owner`, wrong `base_url`) is the top realistic failure, and its failure mode is destruction, not denial.
2. **Provider → GitHub**: token scope is the entire authorisation model. The provider never checks whether the token is broader than the plan needs.
3. **Provider → gh CLI**: an undocumented-in-plan credential source. Silent success here means the operator may not know which identity ran an apply.
4. **Local binary → consumer project**: `dev_overrides` replaces a registry provider with anything on disk. Any process that can write that path controls what every consuming project executes.
5. **Upstream → fork**: the single-touch rule keeps merges small; the same rule means security fixes arrive only as fast as someone merges the release tag.

## Not part of the surface

No listening sockets, no webhooks, no cron, no queue consumers, no database. The plugin is inert until Terraform invokes it, and it opens no port.
