# Layout — `internal/ghclient`

The "new client" implementation, selected with provider argument `legacy_client = false` (default is `true`).

```text
internal/ghclient/
├── source.go          # Source interface: RESTClient/OwnerRESTClient/GraphQLClient/OwnerGraphQLClient
├── token.go           # tokenSource — PAT/OAuth token for REST (go-github v88) and GraphQL (githubv4)
├── app.go             # appSource — GitHub App JWT/installation clients, installation lookup, bounded app-client cache
├── anonymous.go       # anonymousSource — unauthenticated clients, stricter rate limits
├── options.go         # SourceOptions (base URL, GHES flag, UA, cache, retry) → ClientOptions
├── rest.go            # REST client constructors (token / anonymous)
├── graphql.go         # GraphQL client constructors
├── cache.go           # bbolt-backed conditional-request cache (dir 0700, cache.db 0600)
├── throttle.go        # semaphore-based concurrency limiting
├── transport.go       # shared transport assembly (conditional requests, retry)
├── const.go           # constants: concurrency, idle conns, timeouts, UA parts
└── *_test.go          # unit tests for every file above
```

## Rules

- `Source` is the seam: callers never construct clients directly.
- `OwnerRESTClient`/`OwnerGraphQLClient` take an owner; token and anonymous sources ignore it, app source resolves the installation for that owner.
- App source caches clients in a bounded LRU (`appClientCacheSize`) to cap memory.
- Cache path is `filepath.Join(cache_path, "terraform-provider-github", <cacheRef>)`; without `cache_path` a temp dir is used.

## Relationship to `github/transport.go`

`github/transport.go` is the legacy transport stack (`etagTransport` + `RateLimitTransport` + logging + preview headers + `RetryTransport`). `internal/ghclient` replaces it for the new client using the `github-conditional-http-transport` and `go-github-ratelimit` libraries (see `DECISIONS.md`). Both exist in this tree; the switch is `legacy_client`.
