# Architecture — Authentication and Clients

## Credential resolution

`Configure` in `github/provider.go` builds a `Config` and then calls `configureProviderMeta`. Resolution order (legacy client path, the default):

1. **`app_auth` block** (`id`, `installation_id`, `pem_file`) or the `GITHUB_APP_ID` / `GITHUB_APP_INSTALLATION_ID` / `GITHUB_APP_PEM_FILE` env vars. All three are required if any is present; a block with missing fields is an error, not a fallback. The PEM is passed through `strings.ReplaceAll(pem, "\\n", "\n")` so literal `\n` sequences in config/env work. `owner` becomes mandatory.
2. **`token` attribute / `GITHUB_TOKEN`** — PAT or OAuth token.
3. **`gh auth token --hostname <host>`** — `tokenFromGHCLI` shells out to the GitHub CLI (`GH_PATH` overrides the binary). Failure is non-fatal and logged at debug. Host is mapped from the API host (`api.github.com` → `github.com`, `api.<slug>.ghe.com` → `<slug>.ghe.com`).
4. **Anonymous** — unauthenticated clients, public data only, stricter rate limits.

In the **new client path** (`legacy_client = false`) steps 3 does not exist: `ghclient.Source` is chosen from app / token / anonymous, and authenticated use without an explicit `owner` is an error.

Legacy app auth exchanges the JWT for an **installation token** up front (`GenerateOAuthTokenFromApp`, `github/apps.go`), which then becomes the single static token for the whole run. New client app auth instead builds per-owner installation clients lazily (`internal/ghclient/app.go`, bounded cache).

## Client objects

`Owner` (in `config.go`) carries `v3client` (`*github.Client`, REST) and `v4client` (`*githubv4.Client`, GraphQL), plus `name`, `id`, `IsOrganization`. Resources reach them via `m.(*Owner)`.

Base URL handling (`getBaseURL`): the default is `https://api.github.com/`. A `.ghe.com` host means GHEC with the plain REST path; anything else non-default is treated as GHES, which appends `api/v3/` for REST and `api/graphql` for GraphQL.

## When to use which stack

| Concern | Legacy (`legacy_client = true`, default) | New (`legacy_client = false`) |
|---|---|---|
| Client construction | `Config.NewRESTClient` / `NewGraphQLClient` | `ghclient.Source.OwnerRESTClient` / `OwnerGraphQLClient` |
| Rate limiting | `RateLimitTransport` delays + backoff | `go-github-ratelimit` (delay args ignored) |
| Conditional requests | `etagTransport` | `github-conditional-http-transport` + bbolt |
| `gh` CLI fallback | yes | no |
| `parallel_requests` | honoured | always safe (semaphore) |

Deprecated and unused: `insecure` (documented as not applied). `organization` is deprecated in favour of `owner`.
