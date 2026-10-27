# Architecture — Transport and Rate Limiting

## Legacy stack (`github/transport.go`)

`RateLimitedHTTPClient` wraps a base `http.Client` in layers, outermost first:

```
Retry (if max_retries > 0)
└── previewHeaderInjector   # Accept: stone-crop-preview+json
    └── logging             # Terraform log bridge
        └── RateLimit       # read/write delay, abuse-limit backoff, smart lock
            └── etag        # conditional requests + optional cache
                └── default transport
```

- **etagTransport** — stores/uses ETags so repeat reads can return `304` and replay a cached body. With `cache_path` the store persists across runs; without it, caching is per-process.
- **RateLimitTransport** — serializes by request class; write methods (`POST`/`PATCH`/`PUT`/`DELETE`) wait `write_delay_ms` (default 1000), reads wait `read_delay_ms` (default 0). On rate-limit responses it computes the next delay; `parallel_requests` loosens the lock.
- **RetryTransport** — retries `retryable_errors` (default 500/502/503/504) up to `max_retries` (default 3) with `retry_delay_ms` (default 1000).

The body is drained and re-wrapped (`drainBody`) so retries and etag replay can both read it.

## New stack (`internal/ghclient`)

- `github-conditional-http-transport` + its `bbolt` storage: conditional requests with an on-disk store at `<cache_path>/terraform-provider-github/<cacheRef>/cache.db` (dir `0700`, file `0600`).
- `go-github-ratelimit` v2: rate-limit awareness; `read_delay_ms`/`write_delay_ms`/`retryable_errors` are ignored here.
- `throttle.go`: weighted semaphore (`maxConcurrentRequests`) bounds in-flight calls; REST and GraphQL have separate idle-connection/timeout constants in `const.go`.
- Retry via `SourceOptions.RetryMax` / `RetryWaitMin` / `RetryWaitMax` (from `max_retries` and `retry_delay_ms`).

## Practical consequences

- Delay knobs are **request-class based**: every GraphQL call is a `POST`, so it pays `write_delay_ms`. A large read-heavy GraphQL plan is slow under defaults by design (abuse-limit protection).
- `cache_path` is the only persistent write the provider makes; treat its permissions (assumed multi-user host) as a data-exposure surface — it stores response bodies.
- Do not stack both stacks' retry logic; `max_retries` applies to whichever stack is selected.
