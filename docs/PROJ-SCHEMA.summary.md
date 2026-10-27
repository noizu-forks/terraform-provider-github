# Project Schema (summary)

No relational database, no migrations, no ORM. Schema here means: the provider config block, the resource/data-source contracts, Terraform state IDs, and the few on-disk stores.

## Provider config attributes

| Attribute | Default | Env var | Note |
|---|---|---|---|
| `base_url` | `https://api.github.com/` | `GITHUB_BASE_URL` | GHES/GHEC by host |
| `owner` | — | `GITHUB_OWNER` | required for app auth |
| `organization` | — | `GITHUB_ORGANIZATION` | deprecated |
| `token` | — | `GITHUB_TOKEN` | not marked sensitive in schema |
| `app_auth.id` / `.installation_id` / `.pem_file` | — | `GITHUB_APP_*` | pem_file is sensitive |
| `read_delay_ms` / `write_delay_ms` | 0 / 1000 | — | legacy client only |
| `retry_delay_ms` / `max_retries` | 1000 / 3 | — | `retryable_errors` default 500-504 |
| `max_per_page` | 100 | `GITHUB_MAX_PER_PAGE` | |
| `parallel_requests` | false | — | legacy only |
| `insecure` | false | — | deprecated, unused |
| `legacy_client` | true | `GITHUB_LEGACY_CLIENT` | client stack switch |
| `cache_path` | `""` | `GITHUB_CACHE_PATH` | empty ⇒ per-process cache |

Precedence: attribute → `GITHUB_*` env → default. `token` and `app_auth` are not `ConflictsWith` yet (v7 TODO).

## Registries

88 resources / 75 data sources. Largest domains: repository (19/15), actions (18/13), organization (13/16), team (6/1), enterprise (6/1), codespaces (4/6), dependabot (4/4), branch (4/2). Noizu extensions: 0 registered.

## State IDs

Colon-separated composites: `buildID(parts…)` / `parseID2|3|4`; `escapeIDPart` maps `:` → `??`; only the last part may hold an unescaped separator; `buildChecksumID` derives a stable MD5-based ID.

## On-disk data

| Path | Shape |
|---|---|
| `<cache_path>/terraform-provider-github/<ref>/cache.db` | bbolt conditional-request cache (dir 0700, file 0600) |
| `~/.config/opentofu/github-fork.tfrc` | dev_overrides CLI config |
| `~/.local/share/opentofu/github-fork/terraform-provider-github` | built plugin binary |
| `.meta/pointers.yaml` | doc-pointer store (committed) |

## Fixtures

Acceptance resources use the `tf-acc-test-` prefix (`make sweep` cleans up); deploy-key RSA keys are generated at test run time into `test-fixtures/`.

Details: [PROJ-SCHEMA.md](PROJ-SCHEMA.md).
